package owid

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"../configuration"

	"github.com/jinzhu/gorm"
)

// OwidDataSource This is the url to the csv file containing the owid data in the git repository
var OwidDataSource = "https://raw.githubusercontent.com/owid/covid-19-data/master/public/data/owid-covid-data.csv"

//Response Is the structure of the table in the database that will store the owid data
type OwidData struct {
	gorm.Model
	IsoCode               string
	Location              string
	Date                  string
	TotalCases            int64
	NewCases              int64
	TotalDeaths           int64
	NewDeaths             int64
	TotalCasesPerMillion  float64
	NewCasesPerMillion    float64
	TotalDeathsPerMillion float64
	NewDeathsPerMillion   float64
	TotalTests            int64
	NewTests              int64
	TotalTestsPerThousand float64
	NewTestsPerThousand   float64
	TestsUnits            string
}

var myClient = &http.Client{Timeout: 10 * time.Second}

//GetOwidData This function will retrieve the owid data from the git repository
func GetOwidData() ([][]string, error) {

	resp, err := myClient.Get(OwidDataSource)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	data, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err

	}
	//	fmt.Print(data)
	fmt.Print(reflect.TypeOf(data))

	return data, nil
}

func PopulateOwidData(lsRecords [][]string) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user="+configuration.DbUserName+" dbname="+configuration.DbName+" password="+configuration.DbUserPassword+"")
	if err != nil {
		fmt.Println(err.Error())
	}
	if !db.HasTable("OwidData") {
		fmt.Println("creating table OwidData")
		db.CreateTable(&OwidData{})
	}

	for i := range lsRecords {
		if i > 0 {
			//fmt.Println(ls_records[i])

			xOwidRecord := OwidData{
				IsoCode:               lsRecords[i][0],
				Location:              lsRecords[i][1],
				Date:                  lsRecords[i][2],
				TotalCases:            convertStr2Int(lsRecords[i], 3),
				NewCases:              convertStr2Int(lsRecords[i], 4),
				TotalDeaths:           convertStr2Int(lsRecords[i], 5),
				NewDeaths:             convertStr2Int(lsRecords[i], 6),
				TotalCasesPerMillion:  convertStr2Float(lsRecords[i], 7),
				NewCasesPerMillion:    convertStr2Float(lsRecords[i], 8),
				TotalDeathsPerMillion: convertStr2Float(lsRecords[i], 9),
				NewDeathsPerMillion:   convertStr2Float(lsRecords[i], 10),
				TotalTests:            convertStr2Int(lsRecords[i], 11),
				NewTests:              convertStr2Int(lsRecords[i], 12),
				TotalTestsPerThousand: convertStr2Float(lsRecords[i], 13),
				NewTestsPerThousand:   convertStr2Float(lsRecords[i], 14),
				TestsUnits:            lsRecords[i][15]}

			//fmt.Println(xOwidRecord)
			db.NewRecord(xOwidRecord)
			db.Create(&xOwidRecord)

		}
	}

	defer db.Close()

}

func convertStr2Float(lsValues []string, position int) float64 {
	if lsValues[position] != "" {
		value, err := strconv.ParseFloat(lsValues[position], 10)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		return value
	} else {
		return 0.0
	}
}

func convertStr2Int(lsValues []string, position int) int64 {
	if lsValues[position] != "" {
		value, err := strconv.ParseInt(lsValues[position], 10, 64)
		if err != nil {

			tempVal := convertStr2Float(lsValues, position)
			value := int64(tempVal)
			return value
		}
		return value
	} else {
		return 0
	}
}

type Result struct {
	iso_code    string
	date        string
	total_cases int64
}

//IsDataOlderThanDays Will check if the most recent record in the table owid_data is older than x days
func IsDataOlderThanDays(days int64) bool {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user="+configuration.DbUserName+" dbname="+configuration.DbName+" password="+configuration.DbUserPassword+"")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	type Result struct {
		Iso_code string
	}
	rows, err := db.Raw("select iso_code from owid_data where iso_code = ?", "ABW").Rows() // (*sql.Rows, error)
	defer rows.Close()
	for rows.Next() {
		var rq Result
		rows.Scan(rq)
		fmt.Println(rq)
	}

	return true
}
