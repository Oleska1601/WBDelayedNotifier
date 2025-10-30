package usecase

type Usecase struct {
	cache     CacheInterface
	repo      RepoInterface
	publisher PublisherInterface
}

func New(cache CacheInterface, repo RepoInterface, publisher PublisherInterface) *Usecase {
	return &Usecase{
		cache:     cache,
		repo:      repo,
		publisher: publisher,
	}
}
