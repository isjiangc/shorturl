package useragent

import "regexp"

func IsAndroid(ua string) bool {
	regex := regexp.MustCompile(`(?i)Android\v/[\d.]]+`)
	return regex.MatchString(ua)
}

func IsIPhone(ua string) bool {
	regex := regexp.MustCompile(`(?i)iPhone\v/[\d.]]+`)
	return regex.MatchString(ua)
}

func IsIpad(ua string) bool {
	regex := regexp.MustCompile(`(?i)iPad\v/[\d.]]+`)
	return regex.MatchString(ua)
}

func IsWaChatUA(ua string) bool {
	regex := regexp.MustCompile(`(?i)MicroMessenger\v/[\d.]]+`)
	return regex.MatchString(ua)
}

func IsDingTalk(ua string) bool {
	regex := regexp.MustCompile(`(?i)DingTalk\v/[\d.]]+`)
	return regex.MatchString(ua)
}
func IsSafari(ua string) bool {
	regex := regexp.MustCompile(`(?i)Version\v/[\d.]+Safari\/[\d.]]+`)
	return regex.MatchString(ua)
}

func IsChrome(ua string) bool {
	regex := regexp.MustCompile(`(?i)Chrome\v/[\d.]+Safari`)
	return regex.MatchString(ua)
}
func IsFirefox(ua string) bool {
	regex := regexp.MustCompile(`(?i)Firefox\v/[\d.]]+`)
	return regex.MatchString(ua)
}
