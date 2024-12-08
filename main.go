package main

import (
	"fmt"
	"net"
	"time"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	// go server()
	// time.Sleep(5 * time.Second)
	client()

}
func client() {
	inGame := false
	isJoiningRoom := false
	isServerInited := false
	waitingScreen := false
	rl.InitWindow(800, 450, "online pong")
	rl.SetTargetFPS(60)
	var conn *net.UDPConn
	for !rl.WindowShouldClose() {

		if isServerInited {
			serverAddr := &net.UDPAddr{
				IP:   net.IPv4(127, 0, 0, 1),
				Port: 12345,
			}
			conn, _ = net.DialUDP("udp", nil, serverAddr)
			defer conn.Close()
			fmt.Println("conectado ao servidor, digite mensagens")
			isServerInited = false
		}

		if inGame {
			if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
				_, err := conn.Write([]byte("up"))
				if err != nil {
					fmt.Println("erro ao enviar mensagem", err)
					continue
				} else {
					fmt.Println("mensagem enviada")
				}
			}
			if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
				_, err := conn.Write([]byte("down"))
				if err != nil {
					fmt.Println("erro ao enviar mensagem", err)
					continue
				} else {

					fmt.Println("mensagem enviada")
				}
			}

		}
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		if inGame {
			if waitingScreen {
				rl.DrawText("waiting your frind", 400, 225, 20, rl.Black)

			}
			// rl.DrawRectangleRec(player, rl.Red)
			// rl.DrawRectangleRec(enemie, rl.Blue)
			// rl.DrawCircleV(ball.pos, ball.radius, rl.Black)
		} else {
			if raygui.Button(rl.Rectangle{125, 185, 200, 30}, "Init-server") {
				inGame = true
				isServerInited = true
				waitingScreen = true
				go server()
				time.Sleep(2 * time.Second)
			}
			if raygui.Button(rl.Rectangle{445, 185, 200, 30}, "Join a room") && !isJoiningRoom {
				isJoiningRoom = true
			}
			if isJoiningRoom {

			}

		}
		rl.EndDrawing()
	}
	rl.CloseWindow()
}
func server() {
	addr := &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 12345,
		Zone: "",
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("erro ao  iniciar o server", err)
		return
	}
	defer conn.Close()
	fmt.Println("servidor udp escutando na porta 12345")
	buf := make([]byte, 1024)
	player := rl.Rectangle{
		20,
		200,
		20,
		100,
	}
	player2 := rl.Rectangle{
		760,
		200,
		20,
		100,
	}
	ball := struct {
		pos    rl.Vector2
		speed  rl.Vector2
		radius float32
	}{
		pos:    rl.Vector2{400, 225},
		speed:  rl.Vector2{5, 5},
		radius: 10,
	}
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("erro ao ler do client", err)
			continue
		}
		fmt.Printf("recebido do client %v: %s\n", remoteAddr, string(buf[:n]))
		ball.pos.X += ball.speed.X
		ball.pos.Y += ball.speed.Y
		player2.Y = ball.pos.Y
		if rl.CheckCollisionCircleRec(ball.pos, ball.radius, player2) || rl.CheckCollisionCircleRec(ball.pos, ball.radius, player) {
			ball.speed.X = -ball.speed.X
		}
		if rl.CheckCollisionCircleLine(ball.pos, ball.radius, rl.Vector2{0, 0}, rl.Vector2{800, 0}) || rl.CheckCollisionCircleLine(ball.pos, ball.radius, rl.Vector2{0, 450}, rl.Vector2{800, 450}) {
			ball.speed.Y = -ball.speed.Y
		}
		if rl.CheckCollisionCircleLine(ball.pos, ball.radius, rl.Vector2{0, 0}, rl.Vector2{0, 450}) || rl.CheckCollisionCircleLine(ball.pos, ball.radius, rl.Vector2{800, 0}, rl.Vector2{800, 450}) {
			ball.pos = rl.Vector2{400, 250}
		}

	}

}
