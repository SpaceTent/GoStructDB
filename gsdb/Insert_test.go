package gsdb

import (
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	s := NewSuite(t, "sqlite3", ":memory:")

	err := s.CreateTestTables()
	if err != nil {
		t.Fatalf("failed to create test tables for TestInsert: %v", err)
	}

	type Person struct {
		ID        int       `db:"column=id primarykey=yes table=people"`
		Name      string    `db:"column=name"`
		Email     string    `db:"column=email"`
		Active    int       `db:"column=active"`
		CreatedAt time.Time `db:"column=created_at"`
		UpdatedAt time.Time `db:"column=updated_at"`
	}

	testCases := []struct {
		name    string
		person  Person
		wantErr bool
	}{
		{
			name: "successfully insert a Person with all fields filled",
			person: Person{
				Name:      "ZZ",
				Email:     "Top@example.com",
				Active:    1,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
			wantErr: false,
		},
		{
			name: "successfully insert with nullable fields omitted",
			person: Person{
				Name:      "Mr. Robot",
				Active:    1,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
			wantErr: false,
		},
		{
			name: "fail insert with missing non-nullable field",
			person: Person{
				Email:     "failhere@example.com",
				Active:    0,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := s.DB.Insert(tc.person)
			if err != nil {
				t.Fatalf("failed to build insert query: %v", err)
			}

			_, err = s.DB.dbConnection.Exec(query)
			if err != nil {
				t.Fatalf("failed to insert person in tc: %s\n %v", tc.name, err)
			}

			if !tc.wantErr {
				count, err := s.CountRows("people")
				if err != nil {
					t.Fatalf("failed to count rows during test insert: %s\n %v", tc.name, err)
				}

				if count == 0 {
					t.Fatalf("expected at least 1 row during insert test%s\n %v", tc.name, err)
				}
			}
		})
	}

	s.Clear()
	s.TearDown()
}
