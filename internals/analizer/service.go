package analizer

// Service provides analyze operations.
type Service interface {
	Run() (Result, error)
}

type service struct {
	fullPath string
}

// AnalyzeResult is a combined result of repo analisis.
type Result struct {
}

// NewService creates an analize service with the necessary dependencies.
func NewService(path string) Service {
	return &service{fullPath: path}
}

func (s *service) Run() (Result, error) {
	// Let's omit DI pattern for various analyzers here for simplicity
	Scan(s.fullPath)
	return Result{}, nil
}
