package faceguard

import (
	"fmt"

	"github.com/haobird/goutils"
)

func RespAuthorized(client *Client) {
	// photo := "/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAAUDBAQEAwU"
	// str := `{"Num": 1,"PersonInfoList": [{"PersonID": 22,"LastChange": 1602329484,"PersonCode": "5hh","PersonName": "我的陌生人","Remarks": "陌生人的尝试哈哈哈哈哈嘎","TimeTemplateNum": 0,"ImageNum": 1,"ImageList": [{"FaceID": 1,"Name": "1_1.jpg","Size": 3196,"Data": ` + photo + `}]}]}`
	// fmt.Println(len(photo))
	// fmt.Println(len(str))

	// 每 一行 输出一次
	// client.Write([]byte("POST /LAPI/V1.0/PeopleLibraries/3/People HTTP/1.1\r\n"))
	// client.Write([]byte("Content-Type: application/json\r\n"))
	// client.Write([]byte("Connection: close\r\n"))
	// client.Write([]byte("RequestID: 1623999886\r\n"))
	// client.Write([]byte("Content-Length: " + goutils.String(len(str)) + "\r\n\r\n"))

	image := "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"
	photo := ""
	photo = image[0:]
	// photo = image[0:80]
	fmt.Println(photo)
	fmt.Println(len(photo))
	// photo = `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`
	body := `{"Num": 1,"PersonInfoList": [{"PersonID": 22,"LastChange": 1602329484,"PersonCode": "5hh","PersonName": "我的陌生人","Remarks": "陌生人的尝试哈哈哈哈哈嘎","TimeTemplateNum": 0,"ImageNum": 1,"ImageList": [{"FaceID": 1,"Name": "1_1.jpg","Size": ` + goutils.String(len(photo)) + `,"Data": "` + photo + `"}]}]}`

	content := "POST /LAPI/V1.0/PeopleLibraries/3/People HTTP/1.1\r\n" +
		"Content-Type: application/json\r\n" +
		"Connection: close\r\n" +
		"RequestID: 1623999886\r\n" +
		"Content-Length: " + goutils.String(len(body)) + "\r\n\r\n" + body

	fmt.Println("发送消息的长度", len(content))
	client.Write([]byte(content))

	// photo := `/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAAUDBAQEAwUEBAQFBQUGBwwIBwcHBw8LCwkMEQ8SEhEPERETFhwXExQaFRERGCEYGh0dHx8fExciJCIeJBweHx7/2wBDAQUFBQcGBw4ICA4eFBEUHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh7/wAARCABAADIDAREAAhEBAxEB/8QAGQAAAgMBAAAAAAAAAAAAAAAABgcEBQgD/8QANhAAAQMCBQIFAwIDCQAAAAAAAQIDBAURAAYSITETQQciUWFxFDKBCKEVI/BicpGSorHB4fH/xAAbAQACAwEBAQAAAAAAAAAAAAADBAECBQAGB//EADIRAAEDAwIEAgkEAwAAAAAAAAEAAgMEESESMQUTQWFRcQYUIoGRobHB8CMy0eFSkvH/2gAMAwEAAhEDEQA/AErBhJKQQN7emH0wrBqEQkkp3HJtjlynU+kvS30R47SnHnVBDaQAdSibAD84g4C5dM55fquXW5MeZDdjym0FQSpFwq3dJGxHuMc1wcuBuoGXqXKrVAdrEdkNRGFFKy+4kLJAGo+g+L/98TpNiovmyr6hBIsFJAN8SpVWqM5qPlTz6Y5Si+hp+sjIkJZQ3dsJsgHzEE3VvbYgjsNwccVCMcqZIr+YVpRTKbIebKvM6UaW0/KjtijpGt3VSQE+fDvwvgZTZFQmpTUKzYBOhN0ME7eXa/yojjgeqkkxdgbIT5LBT8/5RXXaYYryUSgtxvYJ0li4spaST25tb1BvhKXngh0TtiMeITdK6G1pB0Oe/RZf8QsuZi8O6yqElLzURT31LRa+xSrKR1EdrFJKSk7EEg2NiNZhEjboeFHXHpL+XYs1yoxVz5K1FMeKmyG0DbSUnzAi19+L7XFjiWYOkKLm6piy0CRuPxgllK0v+m6FDfySp52nRZK0y3ApSmgpabBFhuN+b4SnJ1Ibk2pRlLQgRZKWGkmykhF1qPZNyfL35HxgbNPVXjdE0HW0k9M4+mfiqTOFacyxlNl+OsSZi3g11HRY6jdRUU+1uOPxhasn5bdQC3OC8Oj4tXFrhpYBew7WAF/ugKg57zA7LZamym3W33UoJcYB06ja9kaSbel8ZkVbLqAJ+S9lX+ivDxE50TCCATgnp53HyRlmPLsPNlBfodcXEW0u6oLyWFNOtK7nQsk23Tf1v22tswyOabkr5zUxsaLxNItvkOHbIA3ysoZqyLUcpZvdiz0WQ3dTZH2rSbgFJ7p/85vjSY4OFwlwbhQOij0/2xZSnl+mGUlml5ggl09Vx9lYTq2TrSUXH+T9sY/FpOUy4wXWaPf/AAj07NbwDsMp6MLSpao9y5rslKCSNCUmxUo+p/ewwjFO8SaW5JwOwBsSe5z547qHxNLdRx9yeg8vkoWfKLErVIUh4uJcZ/moUkBRuAe3e47fGHKiESssU1wPiUlBUhzNnYN0HSMjCFKh1ehp/jUZCw4tlZ0qUjsUHYHb+u2EjR6CHsyF6mP0m9ZZJS1X6TiLAjoe/VFdKpOl2HIizp6YrOpLkeagLVexGylDUnk7gkHtthxjNiCceK8zV1uJI5WNLnWs5pt8hg+RAI65VP4qZRh5qozkF9goeaBVFl3F2lEfuk7XH/IBxPrckUhAZjxuAFkNwsxysiZoYkusGlzFltZQVIbKkmxIuCOR74fFdT/5j4hXuiHwFnLRm6VD12EmIV+xKFAC/v51f44xPShpNM09/sVqcLI5jvJaToMaV9IkyXAwyT9vC3B2BPYc8b4S4TRzmEB5s0/EjuegQuIVMQkOnJ+QP3K4VOdPjpVMLao8qMySuIp0Fp1scqQfUbbn4tj2MMERAjGWnY9QfArzEs8oJecEbi+COpCraDm1rMVJrdNW8zHqMEBxsoV00OMOedlY3228iv7SVYDGx0NSGFt7fn9pmdwlpi8OtfqiqnTm51JamJUkhxAUoavtPcH4wGaIxPLT0V4JWyxh4KCK14g02PIe68VC4EcHXMcWAkAdwD2/O/tjOkDJnW0AnvlIN4uHzcuNhI8UASPH3KLb7jbdKq7yEqIS4mEgBYB5F1XsffBBwyW37h/qFs6Sl14HyGY/ibRVK+xxbjSibWF21Ef6gnGhWRsfEQ8fnj7kQPc25aVrRx2PNpag7F+v07lrY3UDzuQL9+cJ0FTzLOY+x2J/54pOdjXNOpt+yFc006c/TRCpNMqTCOoFLS86FJCd7lI1E82NuMehpJ2NfqleCew/oLFqoHuj0xMcB5/TJSZztGrWVK7TMywg+24jqRXGHQpCHWbhfSVfgK1qt6FN8Vrg18upjt7beI/AtLhjtdNyZGbX36g/hRvI8RstUvw/k1BuclxFUAESCh1JlBe6XEKQPt02sVHY7EbHClc71oNsPa2KWjoJY2yQt/adj9Uj63UKvmpusVBbrTbVMY+qRTg4dDaN7f3lbElR/FthgDGR09gBkrVo6CKBp5Y/ldqTByrKpUSTJrzzL7zCHHWxpshRSCR+DirppASAE8I2W3VJlWbIouZac/Lacjqjy2nF9RJSQgKF+fa+L1jC+B7RvYpfcLWjWZqRQnEP1WqwIDS07mRISi/xc7/jHieECVtSNIuDuqkXCH6z4+5GhlaKb/Eay6DYmLH0pUfdSyNvgHHsmwOKgMJ3St8SvF2ZnKMmkuZZYpsJLgfS4t4uu6gCBY7AXCjtY4PHFoN7qzW2S9fhKTTKjmFqLrixClMl0KHkJIAFuSTccYuXgODepRA0kXUanM5ip9ZU61R3JDFbpak9H6pKHC2RqBFuDZPB5BI74FI+N+L5BWhT0M7oHVAHsjfPilhKalxZTsV1lIcZWW1C6uQSDh0OaRe6zyLJz10lK1POyk9KXZPTCdfbmx27Af0cAfGdRKU5uhoxtuiai5KylWID1ThOPwFBSh1CsOpKgLkkEcD2O9tuxwEFzMIzJBICQMfVCubKa9lysuUl6XGkrSlKtTQNgCLgEHg2sfzgzHahdXCqjLvcJv7+mLLlEpdfhwct5oodV6iVSwHoDqUFQ6qQu17bg7gA+2AyRkva5vvRWOAaQVBezj9HR6EqkuSjUIbWl8SEAtKPoT9x47Y7kanO17FPwV5hp3RNJ9rpi3v+1kKTKrNlTHpTrcfqPOKcXZHckk4YEbALLNJuV//Z`
	// str := `{"Num": 1,"PersonInfoList": [{"PersonID": 22,"LastChange": 1602329484,"PersonCode": "5hh","PersonName": "我的陌生人","Remarks": "陌生人的尝试哈哈哈哈哈嘎","TimeTemplateNum": 0,"ImageNum": 1,"ImageList": [{"FaceID": 1,"Name": "1_1.jpg","Size": 3196,"Data": ` + photo + `}]}]}`
	// fmt.Println(len(photo))
	// fmt.Println(len(str))
	// length := len(str)
	// for start := 0; start < length; {
	// 	end := 500 + start
	// 	if end > length {
	// 		end = length
	// 	}
	// 	temp := str[start:end]
	// 	fmt.Println(temp)
	// 	start = end
	// 	client.Write([]byte(temp))
	// }
}