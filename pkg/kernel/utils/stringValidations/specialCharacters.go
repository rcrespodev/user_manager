package stringValidations

var specialChars = "$°!@|¬#=%&/()¿?¡~*][}{,;:·+><"

func ContainSpecialChars(input string) (contain bool, chars []string) {
	byteInput := []byte(input)
	specialChars = "$°!@|¬#=%&/()¿?¡~*][}{,;:·+><"
	byteSpecialChars := []byte(specialChars)
	for _, char := range byteInput {
		for _, specialChar := range byteSpecialChars {
			if char == specialChar {
				contain = true
				chars = append(chars, string(char))
				break
			}
		}
	}
	return contain, chars
}
