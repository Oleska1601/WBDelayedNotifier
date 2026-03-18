package service

type Service struct {
	cache     CacheI
	repo      RepoI
	publisher PublisherI
}

func New(cache CacheI, repo RepoI, publisher PublisherI) *Service {
	return &Service{
		cache:     cache,
		repo:      repo,
		publisher: publisher,
	}
}
