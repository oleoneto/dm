package runner_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"regexp"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/oleoneto/dm/pkg/ds"
	"github.com/oleoneto/dm/pkg/migrator"
	"github.com/oleoneto/dm/pkg/runner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TrackerTestSuite struct {
	suite.Suite
	logger  *bytes.Buffer
	options PgOptions
	pg      *embeddedpostgres.EmbeddedPostgres
	db      *sql.DB
}

// =======================================
// MARK: - Suite
// =======================================

// MARK: Before all tests
func (suite *TrackerTestSuite) SetupSuite() {
	suite.options = PgOptions{
		Username:       "user",
		Password:       "pass",
		Name:           "migrations",
		MaxConnections: "20",
		Port:           uint32(54321),
	}

	suite.logger = &bytes.Buffer{}

	config := embeddedpostgres.
		DefaultConfig().
		Username(suite.options.Username).
		Password(suite.options.Password).
		Port(suite.options.Port).
		Database(suite.options.Name).
		Version(embeddedpostgres.V13).
		StartTimeout(20 * time.Second).
		RuntimePath("_temp_pg_data").
		Logger(suite.logger) // suppress all database logs

	suite.pg = embeddedpostgres.NewDatabase(config)
}

// MARK: After all tests
func (suite *TrackerTestSuite) TearDownSuite() {
	// Output database logs
	fmt.Println(suite.logger.String())
}

// MARK: Before each test
func (suite *TrackerTestSuite) SetupTest() {
	suite.pg.Start()

	var err error
	suite.db, err = sql.Open(
		"postgres",
		fmt.Sprintf("host=localhost port=%d user=%s password=%s dbname=%s sslmode=disable", suite.options.Port, suite.options.Username, suite.options.Password, suite.options.Name),
	)
	if err != nil {
		suite.Error(err)
	}

	// fmt.Println("Started and connected to database!")
}

// MARK: After each test
func (suite *TrackerTestSuite) TearDownTest() {
	suite.db.Close()
	suite.pg.Stop()

	// fmt.Println("Closed and shut down database.")
}

// MARK: - Run all tests
func TestTrackerTestSuite(t *testing.T) {
	suite.Run(t, new(TrackerTestSuite))
}

// =======================================
// MARK: - Tests
// =======================================

func (suite *TrackerTestSuite) TestRunnerInitialization() {
	assert.Panics(
		suite.T(),
		func() {
			runner.NewRunner(nil, nil, nil, nil, nil, nil, nil)
		},
	)
}

func (suite *TrackerTestSuite) TestIsTracked() {
	r := runner.NewRunner(
		&runner.TrackerOptions{},
		suite.db,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	ctx := context.TODO()
	isTracked := r.IsTracked(ctx)
	assert.Equal(suite.T(), isTracked, false)

	version := r.Version(ctx)
	assert.Equal(suite.T(), "", version)
}

func (suite *TrackerTestSuite) TestEnableAndDisableTracking() {
	r := runner.NewRunner(
		nil,
		suite.db,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	ctx := context.TODO()

	err := r.StartTracking(ctx)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), r.IsTracked(ctx))
	assert.True(suite.T(), r.IsEmpty(ctx))

	err = r.StopTracking(ctx)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), r.IsTracked(ctx))
}

func (suite *TrackerTestSuite) TestAppliedMigrations() {
	r := runner.NewRunner(
		nil,
		suite.db,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	ctx := context.TODO()

	res, err := r.AppliedMigrations(ctx)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), res)
	assert.Equal(suite.T(), 0, res.Size())
}

func (suite *TrackerTestSuite) TestApplyMigrations() {
	migrations := []migrator.Migration{
		{
			Schema:   3,
			FileName: "20221231054542_create_likes",
			Version:  "20221231054542",
			Name:     "CreateLikes",
			Changes: migrator.Changes{
				Up:   []string{"CREATE TABLE likes (id SERIAL, content_id INT NOT NULL)"},
				Down: []string{"DROP TABLE likes"},
			},
		},
		{
			Schema:   3,
			FileName: "20221231054543_create_comments",
			Version:  "20221231054543",
			Name:     "CreateComments",
			Changes: migrator.Changes{
				Up:   []string{"CREATE TABLE comments (id SERIAL, resource_id INT NOT NULL, content TEXT NOT NULL)"},
				Down: []string{"DROP TABLE comments"},
			},
		},
	}

	r := runner.NewRunner(
		&runner.TrackerOptions{Table: "liquibase"},
		suite.db,
		func(string, *regexp.Regexp) []fs.DirEntry { return []fs.DirEntry{&MockDirEntry{}} }, // file loader
		func([]fs.DirEntry, string, *regexp.Regexp) ([]migrator.Migration, error) { return migrations, nil },
		func(context.Context, ds.Queue[migrator.Migration]) error { return nil }, // validator
		runner.ApplyMigrations,
		nil,
	)

	ctx := context.TODO()

	err := r.StartTracking(ctx)
	assert.NoError(suite.T(), err)

	err = r.Apply(ctx) // TODO: Take loaderFunc as argument
	assert.NoError(suite.T(), err)

	res, err := r.AppliedMigrations(ctx)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), res)
	assert.Equal(suite.T(), 2, res.Size())
}
