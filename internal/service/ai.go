package service

import (
	"blossom/internal/entity"
	"context"
)

const longText = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec a diam lectus. Sed sit amet ipsum mauris. Maecenas congue ligula ac quam viverra nec consectetur ante hendrerit. Donec et mollis dolor. Praesent et diam eget libero egestas mattis sit amet vitae augue."

type AIer interface {
	Ask(ctx context.Context, prompt string) (answer string, err error)
	AskStream(ctx context.Context, prompt string) (delta chan entity.Delta, err error)
}

func (svc *service) Ask(ctx context.Context, prompt string) (answer string, err error) {

	// return "test", nil
	return svc.gpt.Ask(ctx, prompt)
}

func (svc *service) AskStream(ctx context.Context, prompt string) (delta chan entity.Delta, err error) {
	delta = make(chan entity.Delta, 100)

	go func() {
		stream, err := svc.gpt.AskStream(ctx, prompt)
		if err != nil {
			return
		}
		defer stream.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				res, err := stream.Recv()
				// if errors.Is(err, io.EOF) {
				// 	fmt.Println("EOF")
				// 	delta <- entity.Delta{Content: "", Err: err}

				// 	return
				// }
				if err != nil {
					delta <- entity.Delta{Content: "", Err: err}

					return
				}

				delta <- entity.Delta{Content: res.Choices[0].Delta.Content, Err: nil}
			}
		}

	}()

	return delta, nil
}
