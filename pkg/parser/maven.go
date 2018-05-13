package parser

const (
	defaultMavenPackageType = "jar"
)

type MavenPom struct {
			GroupID     string `xml:"groupId"`
			ArtifactID  string `xml:"artifactId"`
			Version     string `xml:"version,omitempty"`
			Packaging   string `xml:"packaging,omitempty"`
			Name        string `xml:"name,omitempty"`
			Description string `xml:"description,omitempty"`
			Parent struct{
				GroupID     string `xml:"groupId"`
				ArtifactID  string `xml:"artifactId"`
				Version     string `xml:"version"`
				RelativePath string `xml:"relativePath,omitempty"`
			}`xml:"parent,omitempty"`
			Modules []string `xml:"modules>module,omitempty"`
			Properties struct {
				JavaVersion string `xml:"java.version,omitempty"`
			} `xml:"properties,omitempty"`
	}

func (m *MavenPom) SetDefaults() {
	pom := MavenPom{}
	pom.Packaging = defaultMavenPackageType
	*m = pom
}

func NewMavenPom() *MavenPom {
	out := MavenPom{}
	out.SetDefaults()
	return &out
}
