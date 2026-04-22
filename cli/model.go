package main

type ResumeData struct {
	Name          string      `yaml:"name"`
	Title         string      `yaml:"title"`
	ProfileImage  string      `yaml:"profile_image"`
	Contact       ContactInfo `yaml:"contact"`
	Sections      []Section   `yaml:"sections"`
	SectionBreaks bool        `yaml:"section_breaks"`
	Theme         string      `yaml:"theme"`
	Font          string      `yaml:"font"`
	HideWatermark bool        // Always true in CLI preview
}

type ContactInfo struct {
	Email    string `yaml:"email"`
	Phone    string `yaml:"phone"`
	Location string `yaml:"location"`
	Github   string `yaml:"github"`
	Linkedin string `yaml:"linkedin"`
	Xing     string `yaml:"xing"`
	Website  string `yaml:"website"`
	Twitter  string `yaml:"twitter"`
}

type Section struct {
	Type      string            `yaml:"type"`
	Title     string            `yaml:"title"`
	Content   string            `yaml:"content"`
	Items     []ExperienceEntry `yaml:"items"`
	Groups    []SkillCategory   `yaml:"groups"`
	Languages []LanguageItem    `yaml:"languages"`
}

type LanguageItem struct {
	Name  string `yaml:"name"`
	Level string `yaml:"level"`
}

type ExperienceEntry struct {
	Company      string   `yaml:"company"`
	Position     string   `yaml:"position"`
	Start        string   `yaml:"start"`
	End          string   `yaml:"end"`
	Description  string   `yaml:"description"`
	Technologies []string `yaml:"technologies"`
	Name         string   `yaml:"name"`
	Issuer       string   `yaml:"issuer"`
	Date         string   `yaml:"date"`
	CredentialID string   `yaml:"credential_id"`
	Institution  string   `yaml:"institution"`
	Degree       string   `yaml:"degree"`
	Field        string   `yaml:"field"`
}

type SkillCategory struct {
	Category string   `yaml:"category"`
	Items    []string `yaml:"items"`
}
