package test

// This test file verifies `UnitSvc.ApplyOps` against the real `CandOrder`
// contract and concurrent multi-user edit semantics.
//
// Covered scenarios:
// 1. `CandOrder` must include every submitted `UnitCre` and `UnitSave` id.
// 2. `CandOrder` must exclude every submitted `UnitDel` id and duplicate ids.
// 3. `UnitCre.LocalId` is resolved to a generated real id before reindexing.
// 4. Under concurrent submissions, the later submitter wins on saved fields,
//    deletions, and final ordering, while unknown units are clustered by index.

import (
	"testing"
	"time"

	"poprako-s/internal/domain/model/aggr"
	"poprako-s/internal/domain/model/query"
	repo_iface "poprako-s/internal/domain/repo"
	repo_mock "poprako-s/internal/domain/repo/mock"
	"poprako-s/internal/domain/svc"
	svc_res "poprako-s/internal/domain/svc/res"
)

// `stubAssignmentRepo` is a minimal assignment repo double for service tests.
type stubAssignmentRepo struct {
	// `assignment` is the value returned by `GetByChapterUserId`.
	assignment *aggr.Assignment

	// `err` is the error returned by `GetByChapterUserId`.
	err repo_iface.RepoErr
}

// `GetById` returns no assignment for unsupported test paths.
func (*stubAssignmentRepo) GetById(string) (*aggr.Assignment, repo_iface.RepoErr) {
	return nil, nil
}

// `GetByChapterUserId` returns the configured assignment lookup result.
func (r *stubAssignmentRepo) GetByChapterUserId(string, string) (*aggr.Assignment, repo_iface.RepoErr) {
	return r.assignment, r.err
}

// `List` returns no assignments for unsupported test paths.
func (*stubAssignmentRepo) List(*query.ListAssignmentOpt) ([]*aggr.Assignment, repo_iface.RepoErr) {
	return nil, nil
}

// `Create` returns no assignment for unsupported test paths.
func (*stubAssignmentRepo) Create(*aggr.AssignmentCre) (*aggr.Assignment, repo_iface.RepoErr) {
	return nil, nil
}

// `Put` performs no action for unsupported test paths.
func (*stubAssignmentRepo) Put(*aggr.AssignmentPut) repo_iface.RepoErr {
	return nil
}

// `Delete` performs no action for unsupported test paths.
func (*stubAssignmentRepo) Delete(string) repo_iface.RepoErr {
	return nil
}

// `DeleteByChapterUserId` performs no action for unsupported test paths.
func (*stubAssignmentRepo) DeleteByChapterUserId(string, string) repo_iface.RepoErr {
	return nil
}

// `DeleteByChapterId` performs no action for unsupported test paths.
func (*stubAssignmentRepo) DeleteByChapterId(string) repo_iface.RepoErr {
	return nil
}

// `stubErrClsf` is a no-op repo error classifier for service tests.
type stubErrClsf struct{}

// `IsNotFound` always returns false.
func (stubErrClsf) IsNotFound(repo_iface.RepoErr) bool {
	return false
}

// `IsTimeout` always returns false.
func (stubErrClsf) IsTimeout(repo_iface.RepoErr) bool {
	return false
}

// `IsUnavailable` always returns false.
func (stubErrClsf) IsUnavailable(repo_iface.RepoErr) bool {
	return false
}

// `mustApplyOpsAccept` runs `ApplyOps` and fails the test if it rejects.
func mustApplyOpsAccept(t *testing.T, repo repo_iface.UnitRepo, diff *aggr.UnitDiff) {
	t.Helper()

	res := svc.UnitSvc{}.ApplyOps(diff, repo)
	if res.IsReject() {
		t.Fatalf("expected accept, got reject code=%d msg=%q", res.Code(), res.Msg())
	}
}

// `mustRejectBadRequest` runs `ApplyOps` and asserts a `BadRequest` rejection.
func mustRejectBadRequest(t *testing.T, pageId string, initUnits []*aggr.Unit, diff *aggr.UnitDiff) {
	t.Helper()

	repo := repo_mock.NewMockUnitRepo(initUnits)
	res := svc.UnitSvc{}.ApplyOps(diff, repo)
	if res.IsAccept() {
		t.Fatalf("expected bad request reject, got accept")
	}
	if res.Code() != svc_res.BadRequest {
		t.Fatalf("expected bad request code, got=%d msg=%q", res.Code(), res.Msg())
	}

	units := repo_mock.ListMockUnitsByPage(repo, pageId)
	if len(units) != len(initUnits) {
		t.Fatalf("reject path should not mutate repo, got=%d want=%d", len(units), len(initUnits))
	}
	for i := range units {
		if units[i].Id != initUnits[i].Id || units[i].Index != initUnits[i].Index {
			t.Fatalf("reject path should keep repo unchanged, got=%+v want=%+v", units[i], initUnits[i])
		}
	}
}

// `mustPageOrder` returns page unit ids sorted by final persisted `Index`.
func mustPageOrder(t *testing.T, repo repo_iface.UnitRepo, pageId string) []string {
	t.Helper()

	units := repo_mock.ListMockUnitsByPage(repo, pageId)
	ids := make([]string, len(units))
	for i := range units {
		ids[i] = units[i].Id
	}

	return ids
}

// `assertOrderEqual` checks exact id sequence equality.
func assertOrderEqual(t *testing.T, got []string, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("order length mismatch, got=%d want=%d, got=%v want=%v", len(got), len(want), got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("order mismatch at %d, got=%v want=%v", i, got, want)
		}
	}
}

// `assertStringPtrEqual` checks pointer equality by content, including `nil`.
func assertStringPtrEqual(t *testing.T, got *string, want *string) {
	t.Helper()

	if got == nil && want == nil {
		return
	}
	if got == nil || want == nil {
		t.Fatalf("string pointer mismatch, got=%v want=%v", got, want)
	}
	if *got != *want {
		t.Fatalf("string value mismatch, got=%q want=%q", *got, *want)
	}
}

// `strPtr` allocates one string pointer for test fixtures.
func strPtr(v string) *string {
	return &v
}

// `TestUnitApplyOps_RejectsCandOrderMissingSubmittedSaveId` verifies that
// `ApplyOps` rejects a diff when `CandOrder` omits any submitted `UnitSave` id.
func TestUnitApplyOps_RejectsCandOrderMissingSubmittedSaveId(t *testing.T) {
	const pageId = "p1"

	initUnits := []*aggr.Unit{{Id: "u0", PageId: pageId, Index: 0}}
	diff := &aggr.UnitDiff{
		PageId: pageId,
		Ops: []aggr.UnitOp{
			&aggr.UnitSave{Id: "u0", PageId: pageId},
		},
		CandOrder: []string{},
	}

	mustRejectBadRequest(t, pageId, initUnits, diff)
}

// `TestUnitApplyOps_RejectsCandOrderContainingDeletedId` verifies that
// `ApplyOps` rejects a diff when `CandOrder` still contains a deleted unit id.
func TestUnitApplyOps_RejectsCandOrderContainingDeletedId(t *testing.T) {
	const pageId = "p1"

	initUnits := []*aggr.Unit{
		{Id: "u0", PageId: pageId, Index: 0},
		{Id: "u1", PageId: pageId, Index: 1},
	}
	diff := &aggr.UnitDiff{
		PageId: pageId,
		Ops: []aggr.UnitOp{
			&aggr.UnitDel{Id: "u1"},
		},
		CandOrder: []string{"u0", "u1"},
	}

	mustRejectBadRequest(t, pageId, initUnits, diff)
}

// `TestUnitApplyOps_RejectsCandOrderContainingDuplicateIds` verifies that
// `ApplyOps` rejects a diff when `CandOrder` repeats any unit id.
func TestUnitApplyOps_RejectsCandOrderContainingDuplicateIds(t *testing.T) {
	const pageId = "p1"

	initUnits := []*aggr.Unit{{Id: "u0", PageId: pageId, Index: 0}}
	diff := &aggr.UnitDiff{
		PageId:    pageId,
		CandOrder: []string{"u0", "u0"},
	}

	mustRejectBadRequest(t, pageId, initUnits, diff)
}

// `TestUnitApplyOps_ResolvesCreateLocalIdToRealId` verifies that
// `CandOrder` may reference a client `LocalId`, and `ApplyOps` remaps it to the
// generated real id before persisting the final order.
func TestUnitApplyOps_ResolvesCreateLocalIdToRealId(t *testing.T) {
	const pageId = "p1"

	repo := repo_mock.NewMockUnitRepo(nil)
	diff := &aggr.UnitDiff{
		PageId: pageId,
		Ops: []aggr.UnitOp{
			&aggr.UnitCre{LocalId: "tmp-1", PageId: pageId},
		},
		CandOrder: []string{"tmp-1"},
	}

	mustApplyOpsAccept(t, repo, diff)

	got := mustPageOrder(t, repo, pageId)
	if len(got) != 1 {
		t.Fatalf("expected exactly one unit, got=%v", got)
	}
	if got[0] == "tmp-1" {
		t.Fatalf("expected generated real id, got temp id: %v", got)
	}
}

// `TestUnitApplyOps_UsesLaterSubmitAsFinalTruthUnderConcurrentEdits` verifies
// the expected concurrent-edit behavior: the later submitter wins on final
// field values, deletions, and order, while units unknown to that submitter are
// clustered near their original index neighbors.
func TestUnitApplyOps_UsesLaterSubmitAsFinalTruthUnderConcurrentEdits(t *testing.T) {
	const pageId = "p1"

	repo := repo_mock.NewMockUnitRepo([]*aggr.Unit{
		{Id: "u0", PageId: pageId, Index: 0, TranslatedText: nil},
		{Id: "u1", PageId: pageId, Index: 1, TranslatedText: nil},
		{Id: "u2", PageId: pageId, Index: 2, TranslatedText: strPtr("base-u2")},
	})

	// Apply client A first. A creates one new unit and saves two existing units.
	diffA := &aggr.UnitDiff{
		PageId: pageId,
		Ops: []aggr.UnitOp{
			&aggr.UnitCre{LocalId: "tmp-a", PageId: pageId, TranslatedText: strPtr("from-a-new")},
			&aggr.UnitSave{Id: "u0", PageId: pageId, TranslatedText: strPtr("from-a-u0")},
			&aggr.UnitSave{Id: "u1", PageId: pageId, TranslatedText: strPtr("from-a-u1")},
		},
		CandOrder: []string{"tmp-a", "u0", "u1", "u2"},
	}
	mustApplyOpsAccept(t, repo, diffA)

	unitsAfterA := repo_mock.ListMockUnitsByPage(repo, pageId)
	if len(unitsAfterA) != 4 {
		t.Fatalf("expected A submit to create one new unit, got=%d units", len(unitsAfterA))
	}
	newUnitId := unitsAfterA[0].Id
	if newUnitId == "tmp-a" || newUnitId == "u0" || newUnitId == "u1" || newUnitId == "u2" {
		t.Fatalf("expected first unit after A to be a generated new id, got=%q", newUnitId)
	}

	// Apply client B afterwards from the same original snapshot. B deletes A's
	// modified `u1`, resets `u0` back to `nil`, and keeps `u2` near the front.
	diffB := &aggr.UnitDiff{
		PageId: pageId,
		Ops: []aggr.UnitOp{
			&aggr.UnitDel{Id: "u1"},
			&aggr.UnitSave{Id: "u0", PageId: pageId, TranslatedText: nil},
			&aggr.UnitSave{Id: "u2", PageId: pageId, TranslatedText: strPtr("from-b-u2")},
		},
		CandOrder: []string{"u2", "u0"},
	}
	mustApplyOpsAccept(t, repo, diffB)

	gotOrder := mustPageOrder(t, repo, pageId)
	assertOrderEqual(t, gotOrder, []string{"u2", newUnitId, "u0"})

	u0 := repo_mock.GetMockUnitById(repo, "u0")
	if u0 == nil {
		t.Fatalf("expected `u0` to exist after concurrent saves")
	}
	assertStringPtrEqual(t, u0.TranslatedText, nil)

	u1 := repo_mock.GetMockUnitById(repo, "u1")
	if u1 != nil {
		t.Fatalf("expected `u1` to be deleted by later submitter, got=%+v", u1)
	}

	u2 := repo_mock.GetMockUnitById(repo, "u2")
	if u2 == nil {
		t.Fatalf("expected `u2` to exist after concurrent saves")
	}
	assertStringPtrEqual(t, u2.TranslatedText, strPtr("from-b-u2"))

	newUnit := repo_mock.GetMockUnitById(repo, newUnitId)
	if newUnit == nil {
		t.Fatalf("expected A-created unit to survive B submit")
	}
	assertStringPtrEqual(t, newUnit.TranslatedText, strPtr("from-a-new"))
}

// `TestUnitCanListPageUnits_AcceptsAnyAssignment` verifies that any existing
// chapter assignment is sufficient to list page units.
func TestUnitCanListPageUnits_AcceptsAnyAssignment(t *testing.T) {
	repo := &stubAssignmentRepo{assignment: &aggr.Assignment{Id: "asg-1", ChapterId: "ch-1", UserId: "u-1"}}

	re := svc.UnitSvc{}.CanListPageUnits("u-1", "ch-1", repo, stubErrClsf{})
	if re.IsReject() {
		t.Fatalf("expected accept, got reject code=%d msg=%q", re.Code(), re.Msg())
	}
}

// `TestUnitCanEditPageUnits_AcceptsTranslatorOrProofreader` verifies that only
// translator or proofreader roles can edit page units.
func TestUnitCanEditPageUnits_AcceptsTranslatorOrProofreader(t *testing.T) {
	now := time.Now()
	repo := &stubAssignmentRepo{assignment: &aggr.Assignment{Id: "asg-1", ChapterId: "ch-1", UserId: "u-1", TimedRoles: aggr.TimedRoles{AssignedTranslatorAt: &now}}}

	re := svc.UnitSvc{}.CanEditPageUnits("u-1", "ch-1", repo, stubErrClsf{})
	if re.IsReject() {
		t.Fatalf("expected accept, got reject code=%d msg=%q", re.Code(), re.Msg())
	}
}

// `TestUnitCanEditPageUnits_RejectsMissingAssignment` verifies that callers
// without any chapter assignment cannot edit page units.
func TestUnitCanEditPageUnits_RejectsMissingAssignment(t *testing.T) {
	repo := &stubAssignmentRepo{}

	re := svc.UnitSvc{}.CanEditPageUnits("u-1", "ch-1", repo, stubErrClsf{})
	if re.IsAccept() {
		t.Fatalf("expected reject, got accept")
	}
	if re.Code() != svc_res.Forbidden {
		t.Fatalf("expected forbidden code, got=%d msg=%q", re.Code(), re.Msg())
	}
}

// `TestUnitCanEditPageUnits_RejectsNonTranslatorProofreader` verifies that
// chapter roles other than translator and proofreader cannot edit page units.
func TestUnitCanEditPageUnits_RejectsNonTranslatorProofreader(t *testing.T) {
	now := time.Now()
	repo := &stubAssignmentRepo{assignment: &aggr.Assignment{Id: "asg-1", ChapterId: "ch-1", UserId: "u-1", TimedRoles: aggr.TimedRoles{AssignedReviewerAt: &now}}}

	re := svc.UnitSvc{}.CanEditPageUnits("u-1", "ch-1", repo, stubErrClsf{})
	if re.IsAccept() {
		t.Fatalf("expected reject, got accept")
	}
	if re.Code() != svc_res.Forbidden {
		t.Fatalf("expected forbidden code, got=%d msg=%q", re.Code(), re.Msg())
	}
}

// `TestMockUnitRepo_CountByPage` verifies that the mock repo counts unit
// states consistently with translated and proofread text presence.
func TestMockUnitRepo_CountByPage(t *testing.T) {
	repo := repo_mock.NewMockUnitRepo([]*aggr.Unit{
		{Id: "u0", PageId: "p1", TranslatedText: strPtr("a")},
		{Id: "u1", PageId: "p1", ProofreadText: strPtr("b")},
		{Id: "u2", PageId: "p1", TranslatedText: strPtr("c"), ProofreadText: strPtr("d")},
		{Id: "u3", PageId: "p2", TranslatedText: strPtr("ignored")},
	})

	total, translated, proofread, err := repo.CountByPage("p1")
	if err != nil {
		t.Fatalf("expected nil error, got=%v", err)
	}
	if total != 3 || translated != 2 || proofread != 2 {
		t.Fatalf("unexpected counts got total=%d translated=%d proofread=%d", total, translated, proofread)
	}
}
