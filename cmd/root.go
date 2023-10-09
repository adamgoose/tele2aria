package cmd

import (
	"log"

	"github.com/adamgoose/tele2aria/lib"
	"github.com/spf13/cobra"
	"github.com/zelenin/go-tdlib/client"
)

var rootCmd = &cobra.Command{
	Use:   "tele2aria",
	Short: "tele2aria is a simple Telegram bot that sends attachments to Aria2",
	RunE: lib.RunE(func(tdlibClient *client.Client) error {
		me, err := tdlibClient.GetMe()
		if err != nil {
			return err
		}

		listener := tdlibClient.GetListener()
		defer listener.Close()

		mids := map[int32]int64{}
		fns := map[int32]string{}

		for update := range listener.Updates {
			if update.GetClass() != client.ClassUpdate {
				continue
			}

			switch update.GetType() {
			case client.TypeUpdateFile:
				u := update.(*client.UpdateFile)
				p := float64(u.File.Local.DownloadedSize) / float64(u.File.ExpectedSize) * 100
				log.Printf("Progress: [%s] %f\n", fns[u.File.Id], p)

				if u.File.Local.IsDownloadingCompleted {
					_, err := tdlibClient.AddMessageReaction(&client.AddMessageReactionRequest{
						ChatId:    me.Id,
						MessageId: mids[u.File.Id],
						ReactionType: &client.ReactionTypeEmoji{
							Emoji: "🎉",
						},
					})
					if err != nil {
						log.Printf("AddMessageReaction error: %s\n", err)
					}

					delete(mids, u.File.Id)
					delete(fns, u.File.Id)
				}

			case client.TypeUpdateNewMessage:
				u := update.(*client.UpdateNewMessage)
				if u.Message.ChatId != me.Id {
					continue
				}

				msg := u.Message
				if msg.Content.MessageContentType() != client.TypeMessageVideo {
					continue
				}

				vm := msg.Content.(*client.MessageVideo)
				f := vm.Video.Video
				mids[f.Id] = msg.Id
				fns[f.Id] = vm.Video.FileName

				_, err := tdlibClient.DownloadFile(&client.DownloadFileRequest{
					FileId:      f.Id,
					Priority:    1,
					Offset:      0,
					Limit:       0,
					Synchronous: false,
				})
				if err != nil {
					log.Printf("DownloadFile error: %s\n", err)
				}

				_, err = tdlibClient.AddMessageReaction(&client.AddMessageReactionRequest{
					ChatId:    me.Id,
					MessageId: msg.Id,
					ReactionType: &client.ReactionTypeEmoji{
						Emoji: "👍",
					},
				})
				if err != nil {
					log.Printf("AddMessageReaction error: %s\n", err)
				}
			}
		}

		return nil
	}),
}

func Execute() error {
	return rootCmd.Execute()
}
