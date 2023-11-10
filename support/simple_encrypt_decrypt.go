package support

// SimpleEncryptDecrypt encrypts or decrypts a byte slice using a simple
// algorithm. In other words:
//
// SimpleEncryptDecrypt(SimpleEncryptDecrypt(data)) == data
//
// This is used, for example, for the broadcast message sent by the robot
// when trying to pair with an app.
func SimpleEncryptDecrypt(data []byte) {
	b := byte(7)
	for i := 0; i < len(data); i++ {
		data[i] = data[i] ^ b
		b = (b + 7) ^ 178
	}
}
