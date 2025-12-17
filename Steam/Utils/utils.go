package Utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html"
	"math/big"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// OpenFormInBrowser creates a temporary HTML file with a self-submitting form
// and opens it in the default web browser. This is used to bypass bot detection
// that might trigger on direct programmatic POST requests.
func OpenFormInBrowser(actionURL string, formData map[string]string) error {
	// 1. Build the HTML content for the form.
	var htmlBuilder strings.Builder
	htmlBuilder.WriteString("<!DOCTYPE html><html><head><title>Redirecting...</title>")
	htmlBuilder.WriteString("<style>body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; } .container { text-align: center; } h1 { color: #333; } p { color: #666; }</style>")
	htmlBuilder.WriteString("</head><body><div class=\"container\">")
	htmlBuilder.WriteString("<h1>Processing your request...</h1>")
	htmlBuilder.WriteString(fmt.Sprintf("<form id=\"redirectForm\" action=\"%s\" method=\"POST\">", html.EscapeString(actionURL)))

	for key, value := range formData {
		htmlBuilder.WriteString(fmt.Sprintf("<input type=\"hidden\" name=\"%s\" value=\"%s\">", html.EscapeString(key), html.EscapeString(value)))
	}

	htmlBuilder.WriteString("<noscript><input type=\"submit\" value=\"Click here to continue\"></noscript>")
	htmlBuilder.WriteString("</form>")
	htmlBuilder.WriteString("<p>You are being automatically redirected. If nothing happens, please enable JavaScript or click the button above.</p>")
	htmlBuilder.WriteString("<script>document.getElementById('redirectForm').submit();</script>")
	htmlBuilder.WriteString("</div></body></html>")

	// 2. Write the HTML to a temporary file.
	tmpFile, err := os.CreateTemp("", "steam_redirect_*.html")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	// Schedule the file for deletion, but also handle immediate errors.
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(htmlBuilder.String()); err != nil {
		tmpFile.Close() // Close before returning
		return fmt.Errorf("failed to write to temporary file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// 3. Open the HTML file in the default browser.
	var cmd *exec.Cmd
	filePath := tmpFile.Name()

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", filePath)
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	case "windows":
		cmd = exec.Command("cmd", "/C", "start", strings.Replace(filePath, "&", "^&", -1))
	default:
		fmt.Printf("Unsupported platform. Please open this file in your browser to continue:\n%s\n", filePath)
		// Keep the file for the user to open manually.
		// We need to prevent the deferred os.Remove from running.
		// A simple way is to just wait for user input.
		fmt.Println("Press Enter to continue after you have opened the file...")
		fmt.Scanln()
		return nil
	}

	fmt.Printf("Opening browser to complete the action. If it does not open, please visit:\n%s\n", filePath)

	// We run this in the background and don't wait for it to complete.
	// This allows the main application to continue, and prevents blocking
	// if the browser process doesn't detach.
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start browser command: %w", err)
	}

	return nil
}

// SafeHexString 生成安全的十六进制随机字符串
// 使用加密安全的随机数生成器，适用于密钥、令牌等安全场景
// 参数：n - 字节长度（最终字符串长度为2*n）
// 返回值：十六进制编码的随机字符串
// 示例：SafeHexString(16) 生成32字符的十六进制字符串
func SafeHexString(n int) string {
	bytes := make([]byte, n)

	// 使用加密安全的随机数填充字节数组
	if _, err := rand.Read(bytes); err != nil {
		panic(err) // 随机数生成失败时终止程序，确保安全性
	}

	// 将字节数组编码为十六进制字符串
	return hex.EncodeToString(bytes)
}

// SafeCustomString 生成自定义字符集的安全随机字符串
// 从指定字符集中随机选择字符组成字符串，适用于密码、验证码等场景
// 参数：
//
//	length - 生成字符串的长度
//	charset - 字符集（例如："0123456789"、"abcdefghijklmnopqrstuvwxyz"）
//
// 返回值：由指定字符集组成的随机字符串
// 示例：SafeCustomString(8, "0123456789") 生成8位数字验证码
func SafeCustomString(length int, charset string) string {
	b := make([]byte, length)

	// 为每个位置从字符集中随机选择字符
	for i := range b {
		// 生成[0, len(charset))范围内的安全随机数
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[idx.Int64()]
	}

	return string(b)
}

func WalletConvert(wallet string) int {
	// 正则表达式匹配数字（整数或小数）
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	// 查找所有匹配的数字
	numbers := re.FindAllString(wallet, -1)
	// 将所有数字拼接起来
	moneyStr := strings.Join(numbers, "")

	// 转换为浮点数
	money, err := strconv.ParseFloat(moneyStr, 64)
	if err != nil {
		// 如果转换失败，返回0
		return 0
	}

	// 乘以100后转换为整数（以分为单位）
	return int(money * 100)
}

func FriendCodeToSteamID64(friendCode uint32) uint64 {
	const base = uint64(76561197960265728)
	return base + uint64(friendCode)
}
