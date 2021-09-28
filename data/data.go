package data

type Page struct {
    Title string
    Body  Data
    Attname string
    Attvalue string
}

type CrData struct {
    Title string
    File bool
}

type Datapage struct {
    Title string
    Body [][]string
    Adata Data
}

type Anspage struct {
    Title string
    Body Data
    Query []string
    Ans string
}

type Data []Attribute

type Attribute struct {
    AttrName string
    AttrValues []string
}
