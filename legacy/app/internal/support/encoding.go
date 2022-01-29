package support

func InPlaceEncodeDecode(data []byte) {
	b := byte(7)
	for i := 0; i < len(data); i++ {
		data[i] = data[i] ^ b
		b = (b + 7) ^ 178
	}
}
