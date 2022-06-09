package mock

import (
	"fmt"

	"github.com/jarcoal/httpmock"
)

func CallAuthSuccess() {
	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("%v%v", "https://localhost:8889", "/$SS$/Services/OAuth/Token"),
		httpmock.NewStringResponder(200, ""))
}
