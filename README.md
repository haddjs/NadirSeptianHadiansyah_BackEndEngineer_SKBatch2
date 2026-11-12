NadirSeptianHadiansyah_BackEndEngineer_SKBatch2/
├── main.go # server, routes
├── models.go # Applicant struct, request/response types
├── rules.go # stages, statuses, transition map, position list
├── validate.go # field validation
├── store.go # in-memory store with a mutex
├── handlers.go # create, list, update-stage
├── errors.go # one error response shape
├── rules_test.go # a few tests on transitions and validation
├── go.mod
└── README.md
