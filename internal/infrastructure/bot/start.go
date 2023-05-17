package bot

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	inviteLinkExpiresIn      = 5 * time.Minute
	MsgAlreadyMember         = "🗿 Ты уже участник этого чата"
	MsgTwitchAlreadyLinked   = "Ваш Twitch аккаунт уже привязан к @%s"
	MsgTelegramAlreadyLinked = "Ваш Telegram аккаунт уже привязан к twitch.tv/%s"
	MsgInviteKeyNotFound     = "🗿 Вас нет в списке приглашенных 📋, попробуйте ещё раз авторизоваться через Twitch"
	MsgInviteLink            = "🤩 Добро пожаловать!\n\n🦦Мы тебя уже заждались!\n\n🔗Ссылка для входа в чат:\n%s"
	MsgInviteIsNotYours      = "🤷‍♂️ Это не ваш ключ приглашения"
	chatStatusMember         = "member"
)

func (b *Bot) CmdStart(ctx context.Context, u tg.Update) error {
	inviteKey := u.Message.CommandArguments()
	if inviteKey == "" {
		return nil
	}

	ownerTwitchID, _, err := b.getOwnerIDsByInviteKey(ctx, u, inviteKey)
	if err != nil {
		return err
	}

	subchatTelegramID, err := b.svc.GetSubchatIDByInviteKey(ctx, inviteKey)
	if err != nil {
		return err
	}

	err = b.checkLinkedAccounts(ctx, u, ownerTwitchID)
	if err != nil {
		return err
	}

	err = b.svc.LinkTelegramToTwitchID(ctx, ownerTwitchID, u.Message.From.ID, u.Message.From.UserName)
	if err != nil {
		return err
	}

	err = b.isAlreadyMember(ctx, u, subchatTelegramID)
	if err != nil {
		return err
	}

	inviteLink, err := b.svc.GetSubchatInviteLinkByTwitchID(ctx, subchatTelegramID, ownerTwitchID)
	if err != nil {
		return err
	}

	if inviteLink != "" {
		err = b.RevokeChatInviteLink(subchatTelegramID, inviteLink)
		if err != nil {
			return err
		}
	}

	inviteLink, err = b.CreateChatInviteLink(subchatTelegramID, inviteLinkExpiresIn)
	if err != nil {
		return err
	}

	err = b.svc.SetSubchatInviteLinkByTwitchID(ctx, subchatTelegramID, ownerTwitchID, inviteLink)
	if err != nil {
		return err
	}

	b.bot.Send(tg.NewMessage(u.Message.From.ID, fmt.Sprintf(MsgInviteLink, inviteLink)))

	return nil
}

func (b *Bot) CreateChatInviteLink(chatID int64, expiresIn time.Duration) (inviteLink string, err error) {
	config := tg.CreateChatInviteLinkConfig{
		ChatConfig:         tg.ChatConfig{ChatID: chatID},
		MemberLimit:        1,
		CreatesJoinRequest: false,
		ExpireDate:         int(time.Now().Add(expiresIn).Unix()),
	}

	res, err := b.bot.Request(config)
	if err != nil {
		b.log.Error(err.Error())
		return "", err
	}

	var resp CreateChatInviteLinkResponse
	err = json.Unmarshal(res.Result, &resp)
	if err != nil {
		b.log.Error(err.Error())
		return "", err
	}

	b.log.Info("CreateChatInviteLinkResponse", resp)

	return resp.InviteLink, nil
}

func (b *Bot) RevokeChatInviteLink(subchatTelegramID int64, inviteLink string) error {
	config := tg.RevokeChatInviteLinkConfig{
		InviteLink: inviteLink,
		ChatConfig: tg.ChatConfig{ChatID: subchatTelegramID},
	}

	res, err := b.bot.Request(config)
	if err != nil {
		b.log.Error(err.Error())
		return err
	}

	var resp RevokeChatInviteLinkResponse
	err = json.Unmarshal(res.Result, &resp)
	if err != nil {
		b.log.Error(err.Error())
		return err
	}

	b.log.Info("RevokeChatInviteLinkResponse", resp)

	return nil
}

func (b *Bot) getOwnerIDsByInviteKey(ctx context.Context, u tg.Update, inviteKey string) (ownerTwitchID int64, ownerTelegramID int64, err error) {
	ownerTwitchID, ownerTelegramID, err = b.svc.GetOwnerIDsByInviteKey(ctx, inviteKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			b.bot.Send(tg.NewMessage(u.Message.From.ID, MsgInviteKeyNotFound))
			return 0, 0, fmt.Errorf("%s", MsgInviteKeyNotFound)
		}
		return 0, 0, err
	}

	if ownerTwitchID == 0 {
		b.bot.Send(tg.NewMessage(u.Message.From.ID, MsgInviteKeyNotFound))
		return 0, 0, fmt.Errorf("%s", MsgInviteKeyNotFound)
	}

	if ownerTelegramID != 0 && ownerTelegramID != u.Message.From.ID {
		b.bot.Send(tg.NewMessage(u.Message.From.ID, MsgInviteIsNotYours))
		return 0, 0, fmt.Errorf("%s", MsgInviteIsNotYours)
	}

	return ownerTwitchID, ownerTelegramID, nil
}

func (b *Bot) isAlreadyMember(ctx context.Context, u tg.Update, subchatTelegramID int64) error {
	chatMember, err := b.bot.GetChatMember(tg.GetChatMemberConfig{ChatConfigWithUser: tg.ChatConfigWithUser{ChatID: subchatTelegramID, UserID: u.Message.From.ID}})
	if err != nil {
		return err
	}

	if chatMember.Status == chatStatusMember {
		b.bot.Send(tg.NewMessage(u.Message.From.ID, MsgAlreadyMember))
		return fmt.Errorf("%s", MsgAlreadyMember)
	}

	return nil
}

func (b *Bot) checkLinkedAccounts(ctx context.Context, u tg.Update, ownerTwitchID int64) error {
	userTelegramID, userTelegramUsername, err := b.svc.GetTelegramByTwitchID(ctx, ownerTwitchID)
	if err != nil {
		return err
	}
	if userTelegramID != 0 && userTelegramID != u.Message.From.ID {
		b.bot.Send(tg.NewMessage(u.Message.From.ID, fmt.Sprintf(MsgTwitchAlreadyLinked, userTelegramUsername)))
		return err
	}

	userTwitchID, userTwitchUsername, err := b.svc.GetTwitchByTelegramUsernameOrID(ctx, u.Message.From.ID, u.Message.From.UserName)
	if err != nil {
		return err
	}
	if userTwitchID != 0 && userTwitchID != ownerTwitchID {
		b.bot.Send(tg.NewMessage(u.Message.From.ID, fmt.Sprintf(MsgTelegramAlreadyLinked, userTwitchUsername)))
		return err
	}

	return nil
}
