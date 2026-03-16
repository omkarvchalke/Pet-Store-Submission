package graph

import (
	"database/sql"
	"errors"
	"time"

	"pet-store/backend/internal/auth"
	"pet-store/backend/internal/service"

	"github.com/graphql-go/graphql"
)

type SchemaDeps struct {
	DB       *sql.DB
	ImageDir string
}

func requireMerchant(p graphql.ResolveParams) (*auth.UserContext, error) {
	user, ok := auth.GetUser(p.Context)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	if user.Role != "merchant" {
		return nil, errors.New("forbidden: merchant credentials required")
	}
	if user.StoreID == nil {
		return nil, errors.New("merchant store missing")
	}
	return user, nil
}

func requireCustomer(p graphql.ResolveParams) (*auth.UserContext, error) {
	user, ok := auth.GetUser(p.Context)
	if !ok {
		return nil, errors.New("unauthorized")
	}
	if user.Role != "customer" {
		return nil, errors.New("forbidden: customer credentials required")
	}
	return user, nil
}

func petMap(
	id string,
	storeID string,
	name string,
	species string,
	age int,
	pictureURL string,
	description string,
	createdAt time.Time,
	status string,
	soldAt sql.NullTime,
) map[string]interface{} {
	item := map[string]interface{}{
		"id":          id,
		"storeId":     storeID,
		"name":        name,
		"species":     species,
		"age":         age,
		"pictureUrl":  pictureURL,
		"description": description,
		"createdAt":   createdAt.Format(time.RFC3339),
		"status":      status,
		"soldAt":      nil,
	}

	if soldAt.Valid {
		item["soldAt"] = soldAt.Time.Format(time.RFC3339)
	}

	return item
}

func NewSchema(deps SchemaDeps) (graphql.Schema, error) {
	merchantService := service.NewMerchantService(deps.DB)
	customerService := service.NewCustomerService(deps.DB)

	petType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Pet",
		Fields: graphql.Fields{
			"id":          &graphql.Field{Type: graphql.String},
			"storeId":     &graphql.Field{Type: graphql.String},
			"name":        &graphql.Field{Type: graphql.String},
			"species":     &graphql.Field{Type: graphql.String},
			"age":         &graphql.Field{Type: graphql.Int},
			"pictureUrl":  &graphql.Field{Type: graphql.String},
			"description": &graphql.Field{Type: graphql.String},
			"createdAt":   &graphql.Field{Type: graphql.String},
			"status":      &graphql.Field{Type: graphql.String},
			"soldAt":      &graphql.Field{Type: graphql.String},
		},
	})

	purchaseResultType := graphql.NewObject(graphql.ObjectConfig{
		Name: "PurchasePetResult",
		Fields: graphql.Fields{
			"success": &graphql.Field{Type: graphql.Boolean},
			"message": &graphql.Field{Type: graphql.String},
		},
	})

	checkoutResultType := graphql.NewObject(graphql.ObjectConfig{
		Name: "CheckoutResult",
		Fields: graphql.Fields{
			"success":             &graphql.Field{Type: graphql.Boolean},
			"message":             &graphql.Field{Type: graphql.String},
			"unavailablePetNames": &graphql.Field{Type: graphql.NewList(graphql.String)},
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"customerAvailablePets": &graphql.Field{
				Type: graphql.NewList(petType),
				Args: graphql.FieldConfigArgument{
					"storeSlug": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if _, err := requireCustomer(p); err != nil {
						return nil, err
					}

					storeSlug := p.Args["storeSlug"].(string)
					pets, err := customerService.GetAvailablePets(storeSlug)
					if err != nil {
						return nil, err
					}

					var result []map[string]interface{}
					for _, pet := range pets {
						result = append(result, petMap(
							pet.ID,
							pet.StoreID,
							pet.Name,
							pet.Species,
							pet.Age,
							pet.PictureURL,
							pet.Description,
							pet.CreatedAt,
							pet.Status,
							pet.SoldAt,
						))
					}

					return result, nil
				},
			},

			"merchantUnsoldPets": &graphql.Field{
				Type: graphql.NewList(petType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := requireMerchant(p)
					if err != nil {
						return nil, err
					}

					pets, err := merchantService.GetUnsoldPets(*user.StoreID)
					if err != nil {
						return nil, err
					}

					var result []map[string]interface{}
					for _, pet := range pets {
						result = append(result, petMap(
							pet.ID,
							pet.StoreID,
							pet.Name,
							pet.Species,
							pet.Age,
							pet.PictureURL,
							pet.Description,
							pet.CreatedAt,
							pet.Status,
							pet.SoldAt,
						))
					}

					return result, nil
				},
			},

			"merchantSoldPets": &graphql.Field{
				Type: graphql.NewList(petType),
				Args: graphql.FieldConfigArgument{
					"start": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"end":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := requireMerchant(p)
					if err != nil {
						return nil, err
					}

					start, err := time.Parse(time.RFC3339, p.Args["start"].(string))
					if err != nil {
						return nil, errors.New("invalid start date format, use RFC3339")
					}

					end, err := time.Parse(time.RFC3339, p.Args["end"].(string))
					if err != nil {
						return nil, errors.New("invalid end date format, use RFC3339")
					}

					pets, err := merchantService.GetSoldPets(*user.StoreID, start, end)
					if err != nil {
						return nil, err
					}

					var result []map[string]interface{}
					for _, pet := range pets {
						result = append(result, petMap(
							pet.ID,
							pet.StoreID,
							pet.Name,
							pet.Species,
							pet.Age,
							pet.PictureURL,
							pet.Description,
							pet.CreatedAt,
							pet.Status,
							pet.SoldAt,
						))
					}

					return result, nil
				},
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"merchantCreatePet": &graphql.Field{
				Type: petType,
				Args: graphql.FieldConfigArgument{
					"name":          &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"species":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"age":           &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"pictureBase64": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"description":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := requireMerchant(p)
					if err != nil {
						return nil, err
					}

					name := p.Args["name"].(string)
					species := p.Args["species"].(string)
					age := p.Args["age"].(int)
					pictureBase64 := p.Args["pictureBase64"].(string)
					description := p.Args["description"].(string)

					pet, err := merchantService.CreatePet(
						p.Context,
						*user.StoreID,
						name,
						species,
						age,
						pictureBase64,
						description,
						deps.ImageDir,
					)
					if err != nil {
						return nil, err
					}

					return petMap(
						pet.ID,
						pet.StoreID,
						pet.Name,
						pet.Species,
						pet.Age,
						pet.PictureURL,
						pet.Description,
						pet.CreatedAt,
						pet.Status,
						pet.SoldAt,
					), nil
				},
			},

			"merchantRemovePet": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"petId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := requireMerchant(p)
					if err != nil {
						return nil, err
					}

					petID := p.Args["petId"].(string)
					if err := merchantService.RemovePetIfAvailable(p.Context, petID, *user.StoreID); err != nil {
						return nil, err
					}

					return true, nil
				},
			},

			"customerPurchasePet": &graphql.Field{
				Type: purchaseResultType,
				Args: graphql.FieldConfigArgument{
					"petId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := requireCustomer(p)
					if err != nil {
						return nil, err
					}

					petID := p.Args["petId"].(string)
					success, message, err := customerService.PurchasePet(p.Context, user.UserID, petID)
					if err != nil {
						return nil, err
					}

					return map[string]interface{}{
						"success": success,
						"message": message,
					}, nil
				},
			},

			"customerCheckoutCart": &graphql.Field{
				Type: checkoutResultType,
				Args: graphql.FieldConfigArgument{
					"petIds": &graphql.ArgumentConfig{Type: graphql.NewList(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					user, err := requireCustomer(p)
					if err != nil {
						return nil, err
					}

					rawPetIDs, ok := p.Args["petIds"].([]interface{})
					if !ok || len(rawPetIDs) == 0 {
						return map[string]interface{}{
							"success":             false,
							"message":             "Your cart is empty.",
							"unavailablePetNames": []string{},
						}, nil
					}

					var petIDs []string
					for _, v := range rawPetIDs {
						petIDs = append(petIDs, v.(string))
					}

					success, message, unavailable, err := customerService.CheckoutCart(p.Context, user.UserID, petIDs)
					if err != nil {
						return nil, err
					}

					return map[string]interface{}{
						"success":             success,
						"message":             message,
						"unavailablePetNames": unavailable,
					}, nil
				},
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
}
