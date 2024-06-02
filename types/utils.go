package types

func BuildChannels(sIn, pIn, sOut, pOut Chan) (ReadChannels, WriteChannels) {
	ins := ReadChannels{PUBLIC: MakeReadOnly(pIn), SECRET: MakeReadOnly(sIn)}
	outs := WriteChannels{PUBLIC: MakeWriteOnly(pOut), SECRET: MakeWriteOnly(sOut)}

	return ins, outs
}
