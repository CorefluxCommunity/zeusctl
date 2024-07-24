package yubikey

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/CorefluxCommunity/zeusctl/pkg/utils"
	"github.com/CorefluxCommunity/zeusctl/pkg/yubikeyscard"
)

// ListYubiKeys will output the basic details of connected YubiKeys.
func ListYubiKeys() error {
	// connect YubiKey smart card interface, disconnect on return
	yks := new(yubikeyscard.YubiKeys)
	if err := yks.Connect(); err != nil {
		return err
	}

	defer yks.Disconnect()

	for i, yk := range yks.YubiKeys {
		ard := yk.AppRelatedData
		crd := yk.CardRelatedData

		utils.PrintHeader(fmt.Sprint(i+1, ": ", yk.ReaderLabel))
		utils.PrintKV("Manufacturer", "Yubico")
		utils.PrintKV("Serial number", fmt.Sprintf("%x", ard.AID.Serial))

		if crd.Name != nil {
			utils.PrintKV("Name of cardholder", strings.Replace(fmt.Sprintf("%s", crd.Name), "<<", " ", -1))
		}

		utils.PrintKV("Signature key", fmt.Sprintf("rsa%d/%s",
			binary.BigEndian.Uint16(ard.AlgoAttrSign.RSAModLen[:]),
			utils.FmtFingerprintTerse(ard.Fingerprints.Sign)))
		utils.PrintKV("Encryption key", fmt.Sprintf("rsa%d/%s",
			binary.BigEndian.Uint16(ard.AlgoAttrEnc.RSAModLen[:]),
			utils.FmtFingerprintTerse(ard.Fingerprints.Enc)))
		utils.PrintKV("Authentication key", fmt.Sprintf("rsa%d/%s",
			binary.BigEndian.Uint16(ard.AlgoAttrAuth.RSAModLen[:]),
			utils.FmtFingerprintTerse(ard.Fingerprints.Auth)))

		if i < len(yks.YubiKeys)-1 {
			fmt.Println()
		}
	}

	return nil
}

// ShowYubiKey will search the connected YubiKeys for the specified serial
// number and output the details including smart card and application-related
// data.
func ShowYubiKey(sn string) error {
	// connect YubiKey smart card interface, disconnect on return
	yks := new(yubikeyscard.YubiKeys)
	if err := yks.Connect(); err != nil {
		return err
	}

	defer yks.Disconnect()

	yk := yks.FindBySN(sn)
	if yk == nil {
		return fmt.Errorf("could not locate YubiKey that supports OpenPGP with serial number '%s'", sn)
	}

	ard := yk.AppRelatedData
	crd := yk.CardRelatedData

	utils.PrintHeader("YubiKey Status")

	utils.PrintKV("Reader", yk.ReaderLabel)
	utils.PrintKV("Application ID", fmt.Sprintf("%x%x%x%x%x%x",
		ard.AID.RID, ard.AID.App, ard.AID.Version,
		ard.AID.Manufacturer, ard.AID.Serial, ard.AID.RFU))
	utils.PrintKV("Application type", "OpenPGP")
	utils.PrintKV("Version", fmt.Sprintf("%d.%d", ard.AID.Version[0], ard.AID.Version[1]))
	utils.PrintKV("Manufacturer", "Yubico")
	utils.PrintKV("Serial number", fmt.Sprintf("%x", ard.AID.Serial))
	utils.PrintKV("Name of cardholder", strings.Replace(fmt.Sprintf("%s", crd.Name), "<<", " ", -1))
	utils.PrintKV("Language prefs", string(crd.LanguagePrefs))

	switch crd.Salutation {
	case 0x30:
		utils.PrintKV("Pronoun", "unspecified")
	case 0x31:
		utils.PrintKV("Pronoun", "he")
	case 0x32:
		utils.PrintKV("Pronoun", "she")
	case 0x39:
		utils.PrintKV("Pronoun", "they")
	}

	utils.PrintKV("Max. PIN lengths", fmt.Sprintf("%d %d %d",
		ard.PWStatus.PW1MaxLenFmt,
		ard.PWStatus.PW1MaxLenRC,
		ard.PWStatus.PW3MaxLenFmt))
	utils.PrintKV("PIN retry counter", fmt.Sprintf("%d %d %d",
		ard.PWStatus.PW1RetryCtr,
		ard.PWStatus.PW1RCRetryCtr,
		ard.PWStatus.PW3RetryCtr))

	utils.PrintKV("Signature key", utils.FmtFingerprintTerse(ard.Fingerprints.Sign))
	utils.PrintKV("    algorithm", fmt.Sprintf("rsa%d",
		binary.BigEndian.Uint16(ard.AlgoAttrSign.RSAModLen[:])))
	signGenDate := int64(binary.BigEndian.Uint32(ard.KeyGenDates.Sign[:]))
	utils.PrintKV("    created", time.Unix(signGenDate, 0).String())

	utils.PrintKV("Encryption key", utils.FmtFingerprintTerse(ard.Fingerprints.Enc))
	utils.PrintKV("    algorithm", fmt.Sprintf("rsa%d",
		binary.BigEndian.Uint16(ard.AlgoAttrEnc.RSAModLen[:])))
	encGenDate := int64(binary.BigEndian.Uint32(ard.KeyGenDates.Enc[:]))
	utils.PrintKV("    created", time.Unix(encGenDate, 0).String())

	utils.PrintKV("Authentication key", utils.FmtFingerprintTerse(ard.Fingerprints.Auth))
	utils.PrintKV("    algorithm", fmt.Sprintf("rsa%d",
		binary.BigEndian.Uint16(ard.AlgoAttrAuth.RSAModLen[:])))
	authGenDate := int64(binary.BigEndian.Uint32(ard.KeyGenDates.Auth[:]))
	utils.PrintKV("    created", time.Unix(authGenDate, 0).String())

	return nil
}
