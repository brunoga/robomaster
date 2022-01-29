package dsp

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git.bug-br.org.br/bga/robomasters1/dsp/internal"
	"github.com/google/uuid"
)

var (
	// Extracted from DJI's RoboMaster S1 app.
	dspKey  = []byte("TRoP4GWuc30k6WUp")
	dspIv   = []byte("bP3crVEO6wABzOc0")
	dspMKey = "wwxnMmF8"
)

// File is the representation of a RoboMaster S1 program file (.dsp). It can be
// used to create, read, modify or write them.
type File struct {
	fileName string // Base file name (no guid).
	dji      internal.Dji
}

// New creates a new File instance with the given creator and title. Returns a
// pointer to a File instance and a nil error on success or nil and a non-nil
// error on failure.
func New(creator, title string) (*File, error) {
	return NewWithPythonCode(creator, title, "")
}

// NewWithPythonCode creates a new File instance with the given creator, title
// and pythonCode. Returns a pointer to a File instance and a nil error on
// success or nil and a non-nil error on failure.
func NewWithPythonCode(creator, title, pythonCode string) (*File, error) {
	trimmedCreator := strings.TrimSpace(creator)
	if len(trimmedCreator) == 0 {
		return nil, fmt.Errorf("creator can not be empty")
	}

	trimmedTitle := strings.TrimSpace(title)
	if len(trimmedTitle) == 0 {
		return nil, fmt.Errorf("title can not be empty")
	}

	now := time.Now()

	f := &File{
		strings.ReplaceAll(title, " ", "-"),
		internal.Dji{
			Attribute: internal.Attribute{
				Creator:                   trimmedCreator,
				CreationDate:              now.Format("2006/01/02"),
				ModifyTime:                "",
				FirmwareVersionDependency: "00.00.0000",
				Title:                     trimmedTitle,
				Guid:                      computeGuid(),
				CodeType:                  "python",
				AppMinVersion:             "",
				AppMaxVersion:             "",
				Sign:                      "",
			},
			Code: internal.Code{
				PythonCode: internal.Cdata{
					Cdata: pythonCode,
				},
				ScratchDescription: internal.Cdata{
					Cdata: "",
				},
			},
		},
	}

	return f, nil
}

// Load loads a RoboMaster S1 program file (.dsp) from disk. Returns a pointer
// to a File instance and a nil error on success or nil and a non-nil error on
// failure.
func Load(fileName string) (*File, error) {
	xmlData, err := decodeDsp(fileName)
	if err != nil {
		return nil, err
	}

	var f File
	err = xml.Unmarshal(xmlData, &f.dji)
	if err != nil {
		return nil, err
	}

	// Remove extension (if there is one).
	baseFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	// Remove guid from fileName if present.
	guidIndex := strings.Index(strings.ToLower(baseFileName),
		strings.ToLower(f.dji.Attribute.Guid))
	if guidIndex != -1 {
		baseFileName = baseFileName[0:guidIndex]
	}

	f.fileName = baseFileName

	return &f, nil
}

// SetPythonCode associates the given pythonCode with the given File.
func (f *File) SetPythonCode(pythonCode string) {
	f.dji.Code.PythonCode.Cdata = strings.TrimSpace(pythonCode)
}

// PythonCode returns the Python code associated with the given File.
func (f *File) PythonCode() string {
	return f.dji.Code.PythonCode.Cdata
}

// Save serializes and saves the File instance to disk at the given path as an
// encrypted RoboMaster S1 program file (.dsp). Returns a nil error on success
// or a non-nil error on failure.
func (f *File) Save(path string) error {
	// Generate final filename. i.e: /path/filenameguid.dsp
	fileName := filepath.Join(path, f.fileName+f.dji.Attribute.Guid+
		".dsp")
	fd, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()

	// Set modified time.
	now := time.Now()
	f.dji.Attribute.ModifyTime = now.Format("01/02/2006 15:04:05")

	f.computeSignature()

	xmlData, err := xml.Marshal(f.dji)
	if err != nil {
		return err
	}

	dspData, err := encodeDsp(xmlData)
	if err != nil {
		return err
	}

	_, err = fd.Write(dspData)
	if err != nil {
		return err
	}

	return nil
}

func (f *File) computeSignature() {
	md5Source := []byte(fmt.Sprintf("%s%s%s%s%s%s%s%s%s",
		dspMKey,
		f.dji.Attribute.CreationDate,
		f.dji.Attribute.Title,
		f.dji.Attribute.Creator,
		f.dji.Attribute.FirmwareVersionDependency,
		f.dji.Attribute.Guid,
		f.dji.Code.PythonCode.Cdata,
		f.dji.Code.ScratchDescription.Cdata,
		// In the original C# code, this was an entry in an enum. When
		// concatenating it to a string, C# actually uses the name of
		// the enum instead of its value. The enum names use title case
		// while the actual value we keep track of is lowercase.
		strings.Title(f.dji.Attribute.CodeType)))

	sum := md5.Sum(md5Source)

	// MD5 sum is 16 bytes. The actual signature is an interval of the
	// signature after it is converted to string. It goes as far as to
	// having "half a byte" at the start and end (ie, the interval starts
	// and ends at "half byte" boundaries).
	fullSignature := fmt.Sprintf("%x", sum)

	f.dji.Attribute.Sign = fullSignature[7:23]
}

func computeGuid() string {
	u, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	// Return the computed Guid without the dashes.
	return strings.ReplaceAll(u.String(), "-", "")
}

// pkcs7Pad right-pads the given byte slice with 1 to n bytes, where
// n is the block size. The size of the result is x times n, where x
// is at least 1.
func pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, fmt.Errorf("invalid block size")
	}
	if b == nil || len(b) == 0 {
		return nil, fmt.Errorf("invalid pkcs7 data")
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// pkcs7Unpad validates and unpads data from the given bytes slice.
// The returned value will be 1 to n bytes smaller depending on the
// amount of padding, where n is the block size.
func pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, fmt.Errorf("invalid block size")
	}
	if b == nil || len(b) == 0 {
		return nil, fmt.Errorf("invalid pkcs7 data")
	}
	if len(b)%blocksize != 0 {
		return nil, fmt.Errorf("invalid pkcs7 padding")
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, fmt.Errorf("invalid pkcs7 padding")
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, fmt.Errorf("invalid pkcs7 padding")
		}
	}
	return b[:len(b)-n], nil
}

func decodeDsp(fileName string) ([]byte, error) {
	fd, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	decoder := base64.NewDecoder(base64.StdEncoding, fd)

	cipherText, err := ioutil.ReadAll(decoder)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(dspKey)
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCDecrypter(block, dspIv)
	cbc.CryptBlocks(cipherText, cipherText)

	plainText, err := pkcs7Unpad(cipherText, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func encodeDsp(plainText []byte) ([]byte, error) {
	newPlainText, err := pkcs7Pad(plainText, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(dspKey)
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCEncrypter(block, dspIv)

	cipherText := make([]byte, len(newPlainText))
	cbc.CryptBlocks(cipherText, newPlainText)

	base64Enc := base64.StdEncoding

	base64Text := make([]byte, base64Enc.EncodedLen(len(cipherText)))
	base64Enc.Encode(base64Text, cipherText)

	return base64Text, nil
}
