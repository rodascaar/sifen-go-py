package main

import (
	"encoding/xml"
	"fmt"

	"github.com/rodascaar/sifen-go-py/sifen/models"
	"github.com/rodascaar/sifen-go-py/sifen/types"
)

func main() {
	de := models.NewDE("0180000001")

	// Create some sample data
	de.DE.GOpeDE = models.TgOpeDE{
		ITipEmi: types.TTipEmi_Normal,
		DCodSeg: "12345",
	}

	de.DE.GTimb = models.TgTimb{
		ITiDE:   types.TTiDE_FacturaElectronica,
		DNumTim: 12345678,
		DEst:    "001",
		DPunExp: "001",
		DNumDoc: "0000001",
		DFeIniT: "2024-01-01",
	}

	// Just a smoke test for nested structures
	de.DE.GDatGralOpe = models.TdDatGralOpe{
		DFeEmiDE: "2024-01-01T10:00:00",
		GOpeCom: &models.TgOpeCom{
			ITImp:    types.TTImp_IVA,
			CMoneOpe: types.CMondT_PYG,
		},
		GEmis: models.TgEmis{
			DRucEm:  "80000001",
			DNomEmi: "Empresa Test SA",
			CDepEmi: types.TDepartamento_Capital,
			CCiuEmi: 1,
			DTelEmi: "123456",
			DEmailE: "test@test.com",
			GActEcoList: []models.TgActEco{
				{CActEco: "12345", DDesActEco: "Activadad Test"},
			},
		},
	}

	output, err := xml.MarshalIndent(de, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output))
}
