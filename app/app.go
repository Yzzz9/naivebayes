package app

import(
    "strings"
    "io/ioutil"
    "time"
    "math/rand"

    "naivebayes/data"
)

func LoadFileStr(title string) (string, error) { 
    filename := title + "_attr.txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        WriteFile(filename, "")
        return "", nil
    }
    return string(body), err
}

func CreateAttrData(body string) *data.Data {
    d := &data.Data{}
    if body == "" {
        return d
    }
    lines := strings.Split(string(body), "\n")
    for _, line := range lines {
        attr := strings.Split(line, ":")
        values := strings.Split(attr[1], ",")
        *d = append(*d, data.Attribute{attr[0], values})
    }
    return d
}

func LoadFile(title string) (*data.Data, error) {
    body, err := LoadFileStr(title)
    return CreateAttrData(body), err
}

func WriteFile(filename, body string) error {
    return ioutil.WriteFile(filename, []byte(body), 0600)
}

func SaveFile(title string, d *data.Data) error {
    filename := title + "_attr.txt"
    var slice []string
    for _, attr := range *d {
        var s string
        s += attr.AttrName + ":"
        s += strings.Join(attr.AttrValues, ",")
        slice = append(slice, s)
    }
    body := strings.Join(slice, "\n")
    return WriteFile(filename, body)
}

func CreateRndFile(title string, d *data.Data, n int) error {
    src := rand.NewSource(time.Now().UnixNano())
    r := rand.New(src)
    var slice []string
    filename := title + "_data.txt"
    for i:=0; i<n; i++ {
        var temp []string
        for j:=0; j<len(*d); j++ {
            sz := len((*d)[j].AttrValues)
            temp = append(temp, (*d)[j].AttrValues[r.Intn(sz)])
        }
        slice = append(slice, strings.Join(temp, ","))
    }
    body := strings.Join(slice, "\n")
    return ioutil.WriteFile(filename, []byte(body), 0600)
}

func LoadRndFile(title string) (*[][]string, error) {
    filename := title + "_data.txt"
    tData := [][]string{}
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return &tData, err
    }
    lines := strings.Split(string(body), "\n")
    for _, line := range lines {
        tData = append(tData, strings.Split(line, ","))
    }
    return &tData, err
}

func SaveRndFile(title string, tData *[][]string) error {
    filename := title + "_data.txt"
    var slice []string
    for i:=0; i<len((*tData)); i++ {
        slice = append(slice, strings.Join((*tData)[i], ","))
    }
    body := strings.Join(slice, "\n")
    return WriteFile(filename, body)
}

func FindAns(tData *[][]string, qslice *[]string, d *data.Data) string {
    qid := -1
    for i := range *qslice {
        if (*qslice)[i] == "?" {
            qid = i
            break
        }
    }

    qvalues := (*d)[qid].AttrValues
    occ := make([]int, len(qvalues))
    for i := range *tData {
        for j := range qvalues {
            if (*tData)[i][qid] == qvalues[j] {
                occ[j]++
                break
            }
        }
    }
    
    ans := make([]float64, len(qvalues))
    for i := range ans {
        ans[i] = float64(occ[i]) / float64(len(*tData))
    }

    for i := range qvalues {
        for j := range *d {
            count := 0
            if j == qid {
                continue
            }
            for k := range *tData {
                if (*tData)[k][j] == (*qslice)[j] && (*tData)[k][qid] == qvalues[i] {
                    count++
                }
            }
            ans[i] *= float64(count) / float64(occ[i])
        }
    }

    fin := -1
    max := 0.0
    for i := range ans {
        if ans[i] > max {
            fin = i
            max = ans[i]
        }
    }

    return qvalues[fin]
}
