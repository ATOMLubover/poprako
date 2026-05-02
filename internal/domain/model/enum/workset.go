package enum

// `WorksetIncl` is the include selector for workset relation preloads.
type WorksetIncl string

const (
	// `WorksetInclTeam` includes the `team` relation in workset queries.
	WorksetInclTeam WorksetIncl = "team"
)
