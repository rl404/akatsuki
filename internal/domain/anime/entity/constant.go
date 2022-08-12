package entity

// Type is anime type.
type Type string

// Available anime types.
const (
	TypeTV      Type = "TV"
	TypeOVA     Type = "OVA"
	TypeONA     Type = "ONA"
	TypeMovie   Type = "MOVIE"
	TypeSpecial Type = "SPECIAL"
	TypeMusic   Type = "MUSIC"
	TypeUnknown Type = ""
)

// Status is anime airing status.
type Status string

// Available anime airing status.
const (
	StatusFinished  Status = "FINISHED"
	StatusReleasing Status = "RELEASING"
	StatusNotYet    Status = "NOT_YET"
)

// Season is anime season.
type Season string

// Available anime seasons.
const (
	SeasonWinter Season = "WINTER"
	SeasonSpring Season = "SPRING"
	SeasonSummer Season = "SUMMER"
	SeasonFall   Season = "FALL"
)

// Day is broadcast day.
type Day string

// Available broadcast day.
const (
	DayMonday    Day = "MONDAY"
	DayTuesday   Day = "TUESDAY"
	DayWednesday Day = "WEDNESDAY"
	DayThursday  Day = "THURSDAY"
	DayFriday    Day = "FRIDAY"
	DaySaturday  Day = "SATURDAY"
	DaySunday    Day = "SUNDAY"
	DayOther     Day = "OTHER"
)

// Source is anime source.
type Source string

// Available anime source.
const (
	SourceOriginal     Source = "ORIGINAL"
	SourceManga        Source = "MANGA"
	Source4Koma        Source = "4_KOMA_MANGA"
	SourceWebManga     Source = "WEB_MANGA"
	SourceDigitalManga Source = "DIGITAL_MANGA"
	SourceNovel        Source = "NOVEL"
	SourceLightNovel   Source = "LIGHT_NOVEL"
	SourceVisualNovel  Source = "VISUAL_NOVEL"
	SourceGame         Source = "GAME"
	SourceCardGame     Source = "CARD_GAME"
	SourceBook         Source = "BOOK"
	SourcePictureBook  Source = "PICTURE_BOOK"
	SourceRadio        Source = "RADIO"
	SourceMusic        Source = "MUSIC"
	SourceOther        Source = "OTHER"
	SourceWebNovel     Source = "WEB_NOVEL"   // undocumented
	SourceMixedMedia   Source = "MIXED_MEDIA" // undocumented
)

// Rating is anime rating.
type Rating string

// Available anime rating.
const (
	RatingG     Rating = "G"
	RatingPG    Rating = "PG"
	RatingPG13  Rating = "PG_13"
	RatingR     Rating = "R"
	RatingRPlus Rating = "R+"
	RatingRX    Rating = "RX"
)

// Relation is anime relation type.
type Relation string

// Available anime relation.
const (
	RelationSequel             = "SEQUEL"
	RelationPrequel            = "PREQUEL"
	RelationAlternativeSetting = "ALTERNATIVE_SETTING"
	RelationAlternativeVersion = "ALTERNATIVE_VERSION"
	RelationSideStory          = "SIDE_STORY"
	RelationParentStory        = "PARENT_STORY"
	RelationSummary            = "SUMMARY"
	RelationFullStory          = "FULL_STORY"
	RelationSpinOff            = "SPIN_OFF"
	RelationAdaptation         = "ADAPTATION"
	RelationCharacter          = "CHARACTER"
	RelationOther              = "OTHER"
)
