package provider

// TODO:Think a better structure.
//type Company int
//type ConnectionType int
//type Network int
//
//var NumberToNetwork = map[int]string{
//	1: "mainnet",
//	3: "ropsten",
//	4: "rinkeby",
//	5: "goerli",
//}
//
//const (
//	Undefined Company = iota
//	Infura
//	Alchemy
//	QuickNode
//)
//
//const (
//	Unknown ConnectionType = iota
//	HTTP
//	WSS
//)
//
//const (
//	Nochain Network = iota
//	Mainnet
//	Expanse
//	Ropsten
//	Rinkeby
//	Goerli
//)
//
//type Provider struct {
//	Who      Company
//	ConnType ConnectionType
//	Net      Network
//	key      string
//}
//
//type Config struct {
//	Begin  string
//	Middle string
//}
//
//func (p *Provider) BuildFullUrl(config *Config) (string, error) {
//	wordlist := []string{config.Begin, NumberToNetwork[int(p.Net)], config.Middle, p.key}
//	return strings.Join(wordlist, ""), nil
//}
