package Wallet

import (
	"testing"
)

func TestCreateKeyPairFromClassicalSeed(t *testing.T) {

	test := []string{
		"xrb_33fabqi49i68qn8kkq1gy8att18w9eqtswyxf76s89jhrasui4smxdfb9exb",
		"xrb_1ie5pnkhgg5s79xfb4ca3j1wiptwjgpsh1yjiqf6bzp9yd377jfrynb7f8qj",
		"xrb_1hdwktdgxrhhfwa1agw9dp633k3shjg9ri9rdt3h5nrga1w7zaub5cstu6kr",
		"xrb_1h767a5kf5qimw38169e6qjnx9xsefs4938o79b7kc6r9k39b1ds4cnes9tx",
		"xrb_3b47a6d3qzs5xtm63tppooih95xakmgaemtjhae395yngxdmbrjpd55xd5p6",
		"xrb_3h1x9znkwqocmk5f5pr418x7hx5y4b5jrx8sxa5wrci8i5egiqozrnyt6fyj",
		"xrb_1i773usoss3huapikx9fyzper5hwwbju63yaxskiy5yd6hmo9nchehdbwk61",
		"xrb_3dye1txcnnbzz1c1fjpuqc483i3m1w3sqtutqeosm5scyoexkrdif1qiipxt",
		"xrb_3pguhxmopb1adj1qj8ajippx5uqz475u5c1x9c49w79njc36oa61xmm4pome",
		"xrb_1d8a4ighk8eato1315qxwi9xfub4bry6qbj7xmz9gpex187xk8g5rycj98gu",
		"xrb_18hnhmp8xuexhfw71ehx4xekuqaynw6giyq3pykotxfu1oh51wb6dgm36nec",
	}

	for i := uint32(0); i < uint32(len(test)); i++ {
		pk, _, err := RecoverKeyPairFromClassicalSeed("F78C3C0C973D15FF2C520230BCB9450090227A8165369DB1324E7506A9A05BF5", i)
		addr := pk.CreateAddress()

		if err != nil {
			t.Errorf("creation from seed falied with %s", test[i])
		}

		if addr != Address(test[i]).UpdatePrefix() {
			t.Errorf("creation from seed falied, given %s when expect %s", addr, test[i])
		}

	}

}

func TestRecoverKeyPairFromSeed(t *testing.T) {

	test := []string{
		"xrb_3tkbkyihfp3qgpqn43wowgidkax4ihhcqors94znsd5oxk4ghdiwwjmes9ji",
		"xrb_3fp1fkazof5j6w7hjw4w4z91rhru8bsg1hcgzd1hhh71785wh8rn8bunbp34",
		"xrb_3gfbeysyy3y8j3efh445kfkoj5ycyjdpum5o6k3g9h4y1uqb1q8hrgy9s3rn",
		"xrb_149g43mx7jbw8igae368jueiqkgw51waiiauxzi96o9eu4rmc8p3naw4cy3k",
	}

	for i := uint32(0); i < uint32(len(test)); i++ {
		pk, _, err := RecoverKeyPairFromCoinSeed([]byte("F78C3C0C973D15FF2C520230BCB9450090227A8165369DB1324E7506A9A05BF5"), i)
		addr := pk.CreateAddress()

		if err != nil {
			t.Errorf("creation from seed falied with %s", test[i])
		}

		if addr != Address(test[i]).UpdatePrefix() {
			t.Errorf("creation from seed falied, given %s when expect %s", addr, test[i])
		}

	}

}
