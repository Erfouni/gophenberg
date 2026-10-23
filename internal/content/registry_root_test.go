// SPDX-License-Identifier: Apache-2.0

package content_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/gopherium/gophenberg/internal/content"
)

func TestDerivedRouteWordReservedNamesTheTypeAndTheWord(t *testing.T) {
	t.Parallel()

	err := content.DerivedRouteWordReserved("media", "3D Files", "media")

	if !errors.Is(err, content.ErrRouteWordReserved) {
		t.Errorf("error = %v, want it to carry %v", err, content.ErrRouteWordReserved)
	}
	if code, _ := content.CodeOf(err); code != "derived_route_word_reserved" {
		t.Errorf("code = %q, want derived_route_word_reserved", code)
	}
	held, _ := content.DetailsOf(err)
	if held["type"] != "media" || held["name"] != "3D Files" || held["word"] != "media" {
		t.Errorf("details = %v, want the type media, the name 3D Files and the word media named", held)
	}
	if !strings.Contains(err.Error(), `"3D Files" would answer under "media"`) {
		t.Errorf("message = %q, want it to name the type and the word", err.Error())
	}
}

func TestRegistryRefusesANewTypeClaimingTheRoot(t *testing.T) {
	t.Parallel()

	store := newFakeTypeStore()
	registry := content.NewRegistry(store)
	claimingRoot := carType(t)
	claimingRoot.RouteWord, claimingRoot.Default = "", true

	_, err := registry.Create(t.Context(), claimingRoot)

	if !errors.Is(err, content.ErrRootTaken) {
		t.Fatalf("Create() error = %v, want %v", err, content.ErrRootTaken)
	}
	if errors.Is(err, content.ErrRouteWordTaken) {
		t.Errorf("Create() error = %v, want the root conflict rather than the named word conflict", err)
	}
	if code, named := content.CodeOf(err); named {
		t.Errorf("Create() error code = %q, want the bare conflict the registry raises on a stored type", code)
	}
	if len(store.types) != 1 {
		t.Errorf("the store holds %d types, want the claim refused before it was written", len(store.types))
	}
}
