package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
	// "unsafe"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Ball struct {
	Pos    rl.Vector2 `json:"pos"`
	Speed  rl.Vector2 `json:"speed"`
	Radius float32    `json:"radius"`
}
type Player struct {
	ID    int          `json:"id"`
	Rec   rl.Rectangle `json:"rec"`
	Score int          `json:"score"`
}
type Game struct {
	Player1 Player `json:"player"`
	Player2 Player `json:"player"`
	Ball    Ball   `json:"ball"`
}

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
	var serverAddr *net.UDPAddr
	for !rl.WindowShouldClose() {

		if isServerInited {
			serverAddr = &net.UDPAddr{
				IP:   net.IPv4(127, 0, 0, 1),
				Port: 12345,
			}
			conn, _ = net.DialUDP("udp", nil, serverAddr)
			defer conn.Close()
			fmt.Println("conectado ao servidor, digite mensagens", "\n")
			isServerInited = false
		}
		if inGame {

			if rl.IsKeyDown(rl.KeyW) {
				_, err := conn.Write([]byte("up"))
				if err != nil {
					fmt.Println("erro ao enviar mensagem", err, "\n")
				} else {
					fmt.Println("mensagem enviada")
				}
			}
			if rl.IsKeyDown(rl.KeyS) {
				_, err := conn.Write([]byte("down"))
				if err != nil {
					fmt.Println("erro ao enviar mensagem", err, "\n")
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
				waitingScreen = false

			} else {
				conn.SetDeadline(time.Now().Add(500 * time.Millisecond))
				buf := make([]byte, 1024)
				n, _, err := conn.ReadFrom(buf)
				if err != nil {
					fmt.Println("erro brabo:", err, "\n")
				}
				var game Game
				err = json.Unmarshal(buf[:n], &game)
				fmt.Println("RECEBIDO DO SERVIDOR:", game)
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
	fmt.Println("servidor udp escutando na porta 12345", "\n")
	buf := make([]byte, 1024)
	rec1 := rl.Rectangle{
		20,
		200,
		20,
		100,
	}
	rec2 := rl.Rectangle{
		760,
		200,
		20,
		100,
	}
	ball := Ball{
		Pos:    rl.Vector2{400, 225},
		Speed:  rl.Vector2{5, 5},
		Radius: 10,
	}
	player1 := Player{
		ID:    1,
		Rec:   rec1,
		Score: 0,
	}
	player2 := Player{
		ID:    2,
		Rec:   rec2,
		Score: 0,
	}
	game := Game{
		player1,
		player2,
		ball,
	}
	for {
		conn.SetDeadline(time.Now().Add(500 * time.Millisecond))
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("erro ao ler do client ", err, "\n")
		}
		fmt.Printf("recebido do client %v: %s\n", remoteAddr, string(buf[:n]))
		data, err := json.Marshal(game)
		if err != nil {
			fmt.Println("erro ao serializar:", err, "\n")
		}
		_, err = conn.WriteToUDP(data, remoteAddr)
		if err != nil {
			fmt.Println("erro ao enviar dados:", err, "\n")
		}
		fmt.Println("MANDEI OS DADOS PORRA")
		game.Ball.Pos.X += game.Ball.Speed.X
		game.Ball.Pos.Y += game.Ball.Speed.Y
		game.Player2.Rec.Y = game.Ball.Pos.Y
		if rl.CheckCollisionCircleRec(game.Ball.Pos, game.Ball.Radius, rec2) || rl.CheckCollisionCircleRec(game.Ball.Pos, game.Ball.Radius, game.Player1.Rec) {
			game.Ball.Speed.X = -game.Ball.Speed.X
		}
		if rl.CheckCollisionCircleLine(game.Ball.Pos, game.Ball.Radius, rl.Vector2{0, 0}, rl.Vector2{800, 0}) || rl.CheckCollisionCircleLine(game.Ball.Pos, game.Ball.Radius, rl.Vector2{0, 450}, rl.Vector2{800, 450}) {
			game.Ball.Speed.Y = -game.Ball.Speed.Y
		}
		if rl.CheckCollisionCircleLine(game.Ball.Pos, game.Ball.Radius, rl.Vector2{0, 0}, rl.Vector2{0, 450}) || rl.CheckCollisionCircleLine(game.Ball.Pos, game.Ball.Radius, rl.Vector2{800, 0}, rl.Vector2{800, 450}) {
			game.Ball.Pos = rl.Vector2{400, 250}
		}

	}

}
