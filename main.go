package main

import (
	"io/ioutil"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

func main() {
	// ホスト鍵の読み込み（ここではダミーの鍵としています。実際には安全に生成した鍵を使用してください）
	privateBytes, err := ioutil.ReadFile("path/to/host/key")
	if err != nil {
		log.Fatal("Failed to load private key: ", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}

	// SSHサーバの設定
	config := &ssh.ServerConfig{
		NoClientAuth: false,
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// ここでパスワード認証を行います（実際の運用ではもっと厳格な認証が必要です）
			if c.User() == "username" && string(pass) == "password" {
				return nil, nil
			}
			return nil, err
		},
	}

	config.AddHostKey(private)

	// SSHサーバの起動
	listener, err := net.Listen("tcp", "0.0.0.0:22")
	if err != nil {
		log.Fatalf("Failed to listen on 0.0.0.0:22 (%s)", err)
	}

	log.Printf("Listening on 0.0.0.0:22...")
	for {
		nConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection (%s)", err)
			continue
		}

		// クライアントからの接続を処理するゴルーチンを起動
		go func(nConn net.Conn) {
			_, chans, reqs, err := ssh.NewServerConn(nConn, config)
			if err != nil {
				log.Printf("Failed to handshake (%s)", err)
				return
			}

			// 接続が確立されたことをログに出力
			log.Printf("New SSH connection from %s", nConn.RemoteAddr())

			// リクエストを無視（通常はX11転送などの処理が必要）
			go ssh.DiscardRequests(reqs)

			// チャネルリクエストの処理
			for newChannel := range chans {
				// チャネルのタイプをチェック
				if newChannel.ChannelType() != "session" {
					newChannel.Reject(ssh.UnknownChannelType, "unsupported channel type")
					continue
				}

				channel, _, err := newChannel.Accept()
				if err != nil {
					log.Printf("Could not accept channel (%s)", err)
					continue
				}

				// ここでクライアントに対して何らかの応答を行う
				channel.Write([]byte("SSH server says hello!\n"))
				channel.Close()
			}
		}(nConn)
	}
}
