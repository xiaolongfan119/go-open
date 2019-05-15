package hypnus

type BodyBind struct{}

func (b *BodyBind) Bind(body map[string]string, obj interface{}) error {
	mapBody(obj, body)
	return nil
}
