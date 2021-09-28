package metaplex_tokenmeta

import (
	"fmt"
	"github.com/near/borsh-go"
)

type Metadata struct {
	Key [1]byte
	P [32]byte
	P2 [32]byte
	//Grr [64]byte
	Data Data
}

type Data struct {
	Name string
	Symbol string
	URI string
	SellerFeeBasisPoints uint16
}

func (m *Metadata) Decode(in []byte) error {
	err := borsh.Deserialize(m, in)
	if err != nil {
		return fmt.Errorf("unpack: %w", err)
	}


	return nil
}

func (m *Data) Decode(in []byte) error {
	//err := borsh.Deserialize(m, in)
	//if err != nil {
	//	return fmt.Errorf("unpack: %w", err)
	//}

	count := 0
	for {
		i := in[count:]
		count ++
		err := borsh.Deserialize(m, i)
		if err != nil {
			//return fmt.Errorf("unpack: %w", err)
			fmt.Println("err count:", count, len(i), err)
			continue
		}
		fmt.Println("count:", count, len(i),"name:",m.Name)
	}

	return nil
}

