package core

import "context"

func (r *Resolver) Message() MessageResolver {
	return &messageResolver{r}
}

type messageResolver struct{ *Resolver }

func (r *messageResolver) Author(ctx context.Context, obj *Message) (*MessageAuthor, error) {
	var author MessageAuthor

	r.App.db.First(&author, "id = ?", obj.MessageAuthorID)
	return &author, nil
}

func (r *messageResolver) Thread(ctx context.Context, obj *Message) (*Thread, error) {
	var thread Thread

	r.App.db.First(&thread, "id = ?", obj.ThreadID)
	return &thread, nil
}
