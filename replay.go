package goosr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/bnch/uleb128"
	"github.com/itchio/lzma"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type Replay struct {
	Gamemode      int8
	Version       int32
	BeatmapMD5    string
	PlayerName    string
	ReplayMD5     string
	Count300      uint16
	Count100      uint16
	Count50       uint16
	CountGeki     uint16
	CountKatu     uint16
	CountMiss     uint16
	TotalScore    int32
	Combo         int16
	PerfectCombo  bool
	Mods          uint32
	LifeBarGraph  string
	TimeStamp     time.Time
	ReplayData    CompressedReplay
	OnlineScoreID int64
	AdditionalMod float64
}

type CompressedReplay []byte

type DecompressedReplay []ReplayValue

type ReplayValue struct {
	Ms         int64
	X          float32
	Y          float32
	KeyPressed KeyPressed
}

type KeyPressed struct {
	MouseLeft  bool
	MouseRight bool
	K1         bool
	K2         bool
	Smoke      bool
}

// NewReplay returns an empty Replay struct with TimeStamp
func NewReplay() (r *Replay) {
	r = &Replay{
		TimeStamp: time.Now(),
	}
	return
}

// NewReplayFromFile reads a replay file, parses it and returns a Replay pointer
func NewReplayFromFile(path string) (r *Replay, err error) {
	r = &Replay{}
	file, err := os.ReadFile(path)
	if err != nil {
		return
	}
	err = r.Unmarshal(file)
	if err != nil {
		return
	}
	return
}

// Formats Replay struct to string
func (r *Replay) String() (s string) {
	s = fmt.Sprintf(
		"Gamemode: %s\n"+
			"Version: %d\n"+
			"Beatmap Hash: %s\n"+
			"Player Name: %s\n"+
			"Replay Hash: %s\n"+
			"Count 300s: %d\n"+
			"Count Geki: %d\n"+
			"Count 100s: %d\n"+
			"Count Katu: %d\n"+
			"Count 50s: %d\n"+
			"Count Misses: %d\n"+
			"Total Score: %d\n"+
			"Combo: %d\n"+
			"Mods: %d\n"+
			"Perfect Combo: %t\n"+
			"TimeStamp: %s\n"+
			"LifeBarGraph: %s\n"+
			"ReplayData Lenght: %d\n"+
			"OnlineScoreID: %d\n",
		ParseGamemode(r.Gamemode),
		r.Version,
		r.BeatmapMD5,
		r.PlayerName,
		r.ReplayMD5,
		r.Count300,
		r.CountGeki,
		r.Count100,
		r.CountKatu,
		r.Count50,
		r.CountMiss,
		r.TotalScore,
		r.Combo,
		r.Mods,
		r.PerfectCombo,
		r.TimeStamp,
		r.LifeBarGraph,
		len(r.ReplayData),
		r.OnlineScoreID,
	)
	return
}

// Marshal encodes Replay struct to byte slice
func (r *Replay) Marshal() (b []byte, err error) {
	buf := bytes.NewBuffer([]byte{})

	err = writeInt8(buf, r.Gamemode)
	if err != nil {
		return
	}

	err = writeInt32(buf, r.Version)
	if err != nil {
		return
	}

	err = writeString(buf, r.BeatmapMD5)
	if err != nil {
		return
	}

	err = writeString(buf, r.PlayerName)
	if err != nil {
		return
	}

	err = writeString(buf, r.ReplayMD5)
	if err != nil {
		return
	}

	err = writeUInt16(buf, r.Count300)
	if err != nil {
		return
	}

	err = writeUInt16(buf, r.Count100)
	if err != nil {
		return
	}

	err = writeUInt16(buf, r.Count50)
	if err != nil {
		return
	}

	err = writeUInt16(buf, r.CountGeki)
	if err != nil {
		return
	}

	err = writeUInt16(buf, r.CountKatu)
	if err != nil {
		return
	}

	err = writeUInt16(buf, r.CountMiss)
	if err != nil {
		return
	}

	err = writeInt32(buf, r.TotalScore)
	if err != nil {
		return
	}

	err = writeInt16(buf, r.Combo)
	if err != nil {
		return
	}

	err = writeBool(buf, r.PerfectCombo)
	if err != nil {
		return
	}

	err = writeUInt32(buf, r.Mods)
	if err != nil {
		return
	}

	err = writeString(buf, r.LifeBarGraph)
	if err != nil {
		return
	}

	err = writeInt64(buf, ticksFromTime(r.TimeStamp))
	if err != nil {
		return
	}

	err = writeByteArray(buf, r.ReplayData)
	if err != nil {
		return
	}

	err = writeInt64(buf, r.OnlineScoreID)
	if err != nil {
		return
	}

	if r.AdditionalMod != 0 {
		err = writeFloat64(buf, r.AdditionalMod)
		if err != nil {
			return
		}
	}

	b = buf.Bytes()
	return
}

// Unmarshal decode byte array to Replay struct
func (r *Replay) Unmarshal(data []byte) (err error) {
	var ticks int64
	b := bytes.NewBuffer(data)
	r.Gamemode, err = readInt8(b)
	if err != nil {
		return
	}
	r.Version, err = readInt32(b)
	if err != nil {
		return
	}
	r.BeatmapMD5, err = readString(b)
	if err != nil {
		return
	}
	r.PlayerName, err = readString(b)
	if err != nil {
		return
	}
	r.ReplayMD5, err = readString(b)
	if err != nil {
		return
	}
	r.Count300, err = readUInt16(b)
	if err != nil {
		return
	}
	r.Count100, err = readUInt16(b)
	if err != nil {
		return
	}
	r.Count50, err = readUInt16(b)
	if err != nil {
		return
	}
	r.CountGeki, err = readUInt16(b)
	if err != nil {
		return
	}
	r.CountKatu, err = readUInt16(b)
	if err != nil {
		return
	}
	r.CountMiss, err = readUInt16(b)
	if err != nil {
		return
	}
	r.TotalScore, err = readInt32(b)
	if err != nil {
		return
	}
	r.Combo, err = readInt16(b)
	if err != nil {
		return
	}
	r.PerfectCombo, err = readBool(b)
	if err != nil {
		return
	}
	r.Mods, err = readUInt32(b)
	if err != nil {
		return
	}
	r.LifeBarGraph, err = readString(b)
	if err != nil {
		return
	}
	ticks, err = readInt64(b)
	if err != nil {
		return
	}
	r.TimeStamp = timeFromTicks(ticks)

	r.ReplayData, err = readByteArray(b)
	if err != nil {
		return
	}
	r.OnlineScoreID, err = readInt64(b)
	if err != nil {
		return
	}
	if b.Len() != 0 {
		r.AdditionalMod, err = readFloat64(b)
		if err != nil {
			return
		}
	}
	return
}

// WriteToFile writes a replay file to specified path
func (r *Replay) WriteToFile(path string) (err error) {
	data, err := r.Marshal()
	if err != nil {
		return
	}
	err = os.WriteFile(path, data, fs.ModePerm)
	if err != nil {
		return
	}
	return
}

// Decompress a ReplayData to DecompressedReplay slice
func (cr *CompressedReplay) Decompress() (dr DecompressedReplay, err error) {
	buf := bytes.NewBuffer(*cr)
	reader := lzma.NewReader(buf)
	defer reader.Close()
	x, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	sa := strings.Split(string(x), ",")
	for _, o := range sa {
		op := strings.Split(o, "|")

		if len(op) < 4 {
			fmt.Println(len(op))
			continue
		}
		var (
			ms           int64
			x            float64
			y            float64
			keyBitOffset int
		)
		ms, err = strconv.ParseInt(op[0], 10, 64)
		if err != nil {
			return
		}
		x, err = strconv.ParseFloat(op[1], 32)
		if err != nil {
			return
		}
		y, err = strconv.ParseFloat(op[2], 32)
		if err != nil {
			return
		}
		keyBitOffset, err = strconv.Atoi(op[3])
		if err != nil {
			return
		}
		dr = append(dr, ReplayValue{
			ms,
			float32(x),
			float32(y),
			KeyPressed{
				keyBitOffset&MOUSELEFT > 0,
				keyBitOffset&MOUSERIGHT > 0,
				keyBitOffset&K1 > 0,
				keyBitOffset&K2 > 0,
				keyBitOffset&SMOKE > 0,
			},
		})
	}
	return
}

// Compress a DecompressedReplay slice to CompressedReplay byte slice
func (dr *DecompressedReplay) Compress() (cr CompressedReplay) {
	b := bytes.NewBuffer(cr)
	writer := lzma.NewWriter(b)
	dec := []string{}
	for _, o := range *dr {
		data := []string{}
		var kbs int
		if o.KeyPressed.MouseLeft {
			kbs += MOUSELEFT
		}
		if o.KeyPressed.MouseRight {
			kbs += MOUSERIGHT
		}
		if o.KeyPressed.K1 {
			kbs += K1
		}
		if o.KeyPressed.K2 {
			kbs += K2
		}
		if o.KeyPressed.Smoke {
			kbs += SMOKE
		}
		data = append(data, fmt.Sprintf("%d", o.Ms), fmt.Sprintf("%g", o.X), fmt.Sprintf("%g", o.Y), strconv.Itoa(kbs))
		dec = append(dec, strings.Join(data, "|"))

	}
	out := strings.Join(dec, ",")
	_, err := io.WriteString(writer, out)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = writer.Close()
	if err != nil {
		return
	}
	cr = b.Bytes()
	return
}

func writeInt8(b io.Writer, i int8) error {
	return binary.Write(b, binary.LittleEndian, i)
}

func readInt8(b io.Reader) (i int8, err error) {
	err = binary.Read(b, binary.LittleEndian, &i)
	return
}

func writeBool(b io.Writer, i bool) error {
	return binary.Write(b, binary.LittleEndian, i)
}

func readBool(b io.Reader) (i bool, err error) {
	err = binary.Read(b, binary.LittleEndian, &i)
	return
}

func writeInt16(b io.Writer, i int16) error {
	return binary.Write(b, binary.LittleEndian, i)
}

func readInt16(b io.Reader) (i int16, err error) {
	err = binary.Read(b, binary.LittleEndian, &i)
	return
}

func writeUInt16(b io.Writer, i uint16) error {
	return binary.Write(b, binary.LittleEndian, i)
}

func readUInt16(b io.Reader) (i uint16, err error) {
	err = binary.Read(b, binary.LittleEndian, &i)
	return
}

func writeInt32(b io.Writer, i int32) error {
	return binary.Write(b, binary.LittleEndian, i)
}

func readInt32(b io.Reader) (i int32, err error) {
	err = binary.Read(b, binary.LittleEndian, &i)
	return
}

func writeUInt32(b io.Writer, i uint32) error {
	return binary.Write(b, binary.LittleEndian, i)
}

func readUInt32(b io.Reader) (i uint32, err error) {
	err = binary.Read(b, binary.LittleEndian, &i)
	return
}

func writeInt64(b io.Writer, i int64) error {
	return binary.Write(b, binary.LittleEndian, i)
}

func readInt64(b io.Reader) (i int64, err error) {
	err = binary.Read(b, binary.LittleEndian, &i)
	return
}

func writeFloat64(b io.Writer, i float64) error {
	return binary.Write(b, binary.LittleEndian, i)
}

func readFloat64(b io.Reader) (i float64, err error) {
	err = binary.Read(b, binary.LittleEndian, &i)
	return
}

func writeByteArray(b io.Writer, d []byte) error {
	binary.Write(b, binary.LittleEndian, int32(len(d)))
	return binary.Write(b, binary.LittleEndian, d)
}

func readByteArray(b io.Reader) (i []byte, err error) {
	var l int32
	err = binary.Read(b, binary.LittleEndian, &l)
	if err != nil {
		return
	}
	i = make([]byte, l)
	err = binary.Read(b, binary.LittleEndian, &i)
	return
}

func writeString(b io.Writer, s string) error {
	var d []byte
	if s == "" {
		d = []byte{0}
	} else {
		d = []byte{11}
		d = append(d, uleb128.Marshal(len(s))...)
		d = append(d, []byte(s)...)
	}
	return binary.Write(b, binary.LittleEndian, d)
}

func readString(b io.Reader) (s string, err error) {
	var (
		p = make([]uint8, 1)
	)
	_, err = b.Read(p)
	if err != nil {
		return
	}
	if p[0] == 0 {
		s = ""
	} else if p[0] == 11 {
		sBuffer := make([]byte, uleb128.UnmarshalReader(b))
		_, err = b.Read(sBuffer)
		if err != nil {
			return
		}
		s = string(sBuffer)
	}
	return
}

func timeFromTicks(t int64) time.Time {
	return time.Unix((t-epoch)/10000000, 0).UTC()
}

func ticksFromTime(t time.Time) int64 {
	return t.Unix()*10000000 + epoch
}

func ParseGamemode(i int8) string {
	switch i {
	case OSU:
		return "Standard"
	case TAIKO:
		return "Taiko"
	case CTB:
		return "Catch The Beat"
	case MANIA:
		return "Mania"
	default:
		return "Standard"
	}
}
