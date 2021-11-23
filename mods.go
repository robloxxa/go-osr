package goosr

var ModsEnum = []BanchoMod{
	{0, "NM", "NoMod"},
	{1, "NF", "NoFail"},
	{2, "EZ", "Easy"},
	{4, "TD", "TouchDevice"},
	{8, "HD", "Hidden"},
	{16, "HR", "HardRock"},
	{32, "SD", "SuddenDeath"},
	{64, "DT", "DoubleTime"},
	{128, "RX", "Relax"},
	{256, "HT", "HalfTime"},
	{512, "NC", "NightCore"},
	{1024, "FL", "Flashlight"},
	{2048, "AT", "Auto"},
	{4096, "SO", "SpunOut"},
	{8192, "AP", "AutoPilot"},
	{16384, "PF", "Perfect"},
}

type BanchoMod struct {
	Value uint32
	Short string
	Long  string
}

func ParseBitFlags(bits uint32) (mods []BanchoMod) {
	if bits == 0 {
		return
	}
	for _, mod := range ModsEnum {
		if (bits & mod.Value) != 0 {
			mods = append(mods, mod)
		}
	}
	return
}

func ReturnBitFlags(mods []BanchoMod) (i uint32) {
	for _, mod := range mods {
		i += mod.Value
	}
	return
}
