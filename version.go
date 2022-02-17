package gox

var (
    // BuildDate date string of when build was performed filled in by -X compile flag
    BuildDate string

    // GitRepo string of the git repo url when build was performed filled in by -X compile flag
    GitRepo string

    // BuiltBy date string of who performed buildfilled in by -X compile flag
    BuiltBy string

    // CommitDate date string of when commit of the build was performed filled in by -X compile flag
    CommitDate string

    // Branch string of branch in the git repo filled in by -X compile flag
    Branch string

    // LatestCommit date string of when build was performed filled in by -X compile flag
    LatestCommit string

    // Version string of build filled in by -X compile flag
    Version string
)
