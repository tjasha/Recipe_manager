import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import api from '../api/client';
import { toast } from 'react-toastify';

const IngredientAdder = ({ allIngredients, onAdd }: { allIngredients: any[], onAdd: (id: number, qty: number) => void }) => {
    const [selectedId, setSelectedId] = useState<number | ''>('');
    const [quantity, setQuantity] = useState<number | ''>('');

    const handleAdd = () => {
        if (selectedId !== '' && quantity !== '' && quantity > 0) {
            onAdd(selectedId, quantity);
            setSelectedId('');
            setQuantity('');
        }
    };

    return (
        <div className="input-group mb-3">
            <select className="form-select" value={selectedId} onChange={e => setSelectedId(Number(e.target.value))}>
                <option value="">Choose ingredient...</option>
                {allIngredients.map(ing => (
                    <option key={ing.id} value={ing.id}>{ing.ingredient}</option>
                ))}
            </select>
            <input
                type="number"
                className="form-control"
                placeholder="Quantity"
                value={quantity}
                onChange={e => setQuantity(Number(e.target.value))}
            />
            <button className="btn btn-outline-secondary" type="button" onClick={handleAdd}>Add</button>
        </div>
    );
};


export default function EditRecipePage() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();

    // all states
    const [formData, setFormData] = useState<any | null>(null);
    const [allIngredients, setAllIngredients] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    // to load the data
    useEffect(() => {
        const fetchData = async () => {
            try {
                const recipeResponse = await api.get(`/myrecipe/${id}`);
                const recipeData = recipeResponse.data;

                const ingredientsResponse = await api.get('/ingredients');
                const ingredientsData = ingredientsResponse.data;

                setAllIngredients(ingredientsData);

                // formating data
                const formattedData = {
                    id: recipeData.id,
                    title: recipeData.title,
                    description: recipeData.description || '',
                    portion: recipeData.portion,
                    preparationTime: recipeData.preparationTime ?? '',
                    cookingTime: recipeData.cookingTime ?? '',
                    imageURL: recipeData.imageURL || '',
                    published: recipeData.published,
                    nutrition: {
                        calories: recipeData.nutrition?.calories ?? '',
                        fat: recipeData.nutrition?.fat ?? '',
                        protein: recipeData.nutrition?.protein ?? '',
                        carbohydrate: recipeData.nutrition?.carbohydrate ?? '',
                        sodium: recipeData.nutrition?.sodium ?? '',
                        fiber: recipeData.nutrition?.fiber ?? '',
                        sugar: recipeData.nutrition?.sugar ?? '',
                },
                    ingredients: (recipeData.ingredients ?? []).map((ing: any) => ({
                        ingredientId: ing.ingredient.id,
                        quantity: ing.quantity,
                    })),
                    instructions: (recipeData.instructions ?? []).map((inst: any) => ({
                        stepDescription: inst.stepDescription,
                    })),
                };
                setFormData(formattedData);
            } catch (err: any) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };
        fetchData();
    }, [id]);

    //  same handlers as in create
    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value, type } = e.target;
        if (type === 'checkbox') {
            const { checked } = e.target as HTMLInputElement;
            setFormData((prev: any) => ({ ...prev, [name]: checked }));
            return;
        }
        // For number inputs, allow the state to hold an empty string
        const processedValue = (name === 'portion' && (value === '' || parseFloat(value) < 1)) ? 1 : value;
        setFormData((prev: any) => ({ ...prev, [name]: processedValue }));
    };
    const handleNutritionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData((prev: any) => ({ ...prev, nutrition: { ...prev.nutrition, [name]: value } }));
    };
    const handleInstructionChange = (index: number, value: string) => {
        const newInstructions = [...formData.instructions];
        newInstructions[index].stepDescription = value;
        setFormData((prev: any) => ({ ...prev, instructions: newInstructions }));
    };
    const addInstruction = () => {
        setFormData((prev: any) => ({ ...prev, instructions: [...prev.instructions, { stepDescription: '' }] }));
    };
    const removeInstruction = (index: number) => {
        setFormData((prev: any) => ({ ...prev, instructions: prev.instructions.filter((_: any, i: number) => i !== index) }));
    };
    const addIngredient = (ingredientId: number, quantity: number) => {
        if (!ingredientId || quantity <= 0) return;
        setFormData((prev: any) => ({ ...prev, ingredients: [...prev.ingredients, { ingredientId, quantity }] }));
    };
    const removeIngredient = (index: number) => {
        setFormData((prev: any) => ({ ...prev, ingredients: prev.ingredients.filter((_: any, i: number) => i !== index) }));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!formData) return;
        const nutritionPayload = {
            calories: formData.nutrition.calories === '' ? null : parseFloat(formData.nutrition.calories),
            fat: formData.nutrition.fat === '' ? null : parseFloat(formData.nutrition.fat),
            protein: formData.nutrition.protein === '' ? null : parseFloat(formData.nutrition.protein),
            carbohydrate: formData.nutrition.carbohydrate === '' ? null : parseFloat(formData.nutrition.carbohydrate),
            sodium: formData.nutrition.sodium === '' ? null : parseFloat(formData.nutrition.sodium),
            fiber: formData.nutrition.fiber === '' ? null : parseFloat(formData.nutrition.fiber),
            sugar: formData.nutrition.sugar === '' ? null : parseFloat(formData.nutrition.sugar)
        };

        // Check if all nutrition fields are null
        const isNutritionEmpty = Object.values(nutritionPayload).every(v => v === null);

        const payload = {
            ...formData,
            preparationTime: formData.preparationTime === '' ? null : parseFloat(formData.preparationTime),
            cookingTime: formData.cookingTime === '' ? null : parseFloat(formData.cookingTime),
            // If nutrition is empty, send null, otherwise send the cleaned payload
            nutrition: isNutritionEmpty ? null : nutritionPayload,
        };

        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            await api.put(`/myrecipe/${id}`, payload);
            toast.success('Recipe updated successfully.');
            navigate(`/myrecipe/${id}`);
        } catch (err: any) {
            setError(err.response?.data?.message || err.message || 'Failed to update recipe.');
            toast.error(err.response?.data?.message || err.message || 'Failed to update recipe.');
        } finally {
            setLoading(false);
        }
    };

    if (loading) return <div className="container mt-4"><h2>Loading editor...</h2></div>;
    if (error) return <div className="container mt-4"><div className="alert alert-danger">{error}</div></div>;
    if (!formData) return <div className="container mt-4"><h2>Could not load recipe data.</h2></div>;

    return (
        <div className="container mt-4 mb-5">
            <h1>Edit Recipe</h1>
            <form onSubmit={handleSubmit}>
                {/* --- Basic recipe fields --- */}
                <div className="card p-4 mb-4">
                    <h2>Basic Information</h2>
                    <div className="mb-3">
                        <label htmlFor="title" className="form-label">Title*</label>
                        <input type="text" id="title" name="title" className="form-control" value={formData.title}
                               onChange={handleChange} required/>
                    </div>
                    <div className="mb-3">
                        <label htmlFor="description" className="form-label">Description</label>
                        <textarea id="description" name="description" className="form-control" rows={3}
                                  value={formData.description} onChange={handleChange}></textarea>
                    </div>
                    <div className="row">
                        <div className="col-md-4 mb-3">
                            <label htmlFor="portion" className="form-label">Portion (servings)</label>
                            <input type="number" id="portion" name="portion" className="form-control"
                                   value={formData.portion} onChange={handleChange} min="1"/>
                        </div>
                        <div className="col-md-4 mb-3">
                            <label htmlFor="preparationTime" className="form-label">Preparation Time (min)</label>
                            <input type="number" id="preparationTime" name="preparationTime" className="form-control"
                                   value={formData.preparationTime} onChange={handleChange} min="0"/>
                        </div>
                        <div className="col-md-4 mb-3">
                            <label htmlFor="cookingTime" className="form-label">Cooking Time (min)</label>
                            <input type="number" id="cookingTime" name="cookingTime" className="form-control"
                                   value={formData.cookingTime} onChange={handleChange} min="0"/>
                        </div>
                    </div>
                    <div className="mb-3">
                        <label htmlFor="imageURL" className="form-label">Image URL</label>
                        <input type="text" id="imageURL" name="imageURL" className="form-control"
                               value={formData.imageURL} onChange={handleChange}/>
                    </div>
                </div>

                {/* --- Ingredients --- */}
                <div className="card p-4 mb-4">
                    <h2>Ingredients*</h2>
                    <IngredientAdder allIngredients={allIngredients} onAdd={addIngredient} />                    <ul className="list-group">
                        {formData.ingredients.map((ing, index) => {
                            const ingredientDetails = allIngredients.find(i => i.id === ing.ingredientId);
                            return (
                                <li key={index}
                                    className="list-group-item d-flex justify-content-between align-items-center">
                                    {ingredientDetails ? `${ing.quantity} ${ingredientDetails.unit || ''} ${ingredientDetails.ingredient}` : 'Loading...'}
                                    <button type="button" className="btn btn-sm btn-outline-danger"
                                            onClick={() => removeIngredient(index)}>Remove
                                    </button>
                                </li>
                            );
                        })}
                    </ul>
                </div>

                {/* --- Instructions --- */}
                <div className="card p-4 mb-4">
                    <h2>Instructions</h2>
                    <p/>
                    {formData.instructions.map((inst, index) => (
                        <div key={index} className="input-group mb-2">
                            <span className="input-group-text">{index + 1}.</span>
                            <textarea
                                className="form-control"
                                rows={2}
                                value={inst.stepDescription}
                                onChange={e => handleInstructionChange(index, e.target.value)}
                            ></textarea>
                            <button type="button" className="btn btn-outline-danger"
                                    onClick={() => removeInstruction(index)}>X
                            </button>
                        </div>
                    ))}
                    <button type="button" className="btn btn-outline-secondary mt-2" onClick={addInstruction}>+ Add
                        Step
                    </button>
                </div>

                {/* --- Nutrition Facts --- */}
                <div className="card p-4 mb-4">
                    <h2>Nutrition Facts (per serving)</h2>
                    <p></p>
                    <div className="row">
                        <div className="col-md-3 mb-3"><label className="form-label">Calories (kcal)</label>
                            <input type="number" name="calories" className="form-control"
                                   value={formData.nutrition.calories}
                                   onChange={handleNutritionChange} min="0"/>
                        </div>
                        <div className="col-md-3 mb-3"><label className="form-label">Protein (g)</label>
                            <input type="number" name="protein" className="form-control"
                                   value={formData.nutrition.protein}
                                   onChange={handleNutritionChange} min="0"/>
                        </div>
                        <div className="col-md-3 mb-3"><label className="form-label">Carbohydrate (g)</label>
                            <input type="number" name="carbohydrate" className="form-control"
                                   value={formData.nutrition.carbohydrate}
                                   onChange={handleNutritionChange} min="0"/>
                        </div>
                        <div className="col-md-3 mb-3"><label className="form-label">Fiber (g)</label>
                            <input type="number" name="fiber" className="form-control" value={formData.nutrition.fiber}
                                   onChange={handleNutritionChange} min="0"/>
                        </div>
                        <div className="col-md-3 mb-3"><label className="form-label">Fat (g)</label>
                            <input type="number" name="fat" className="form-control" value={formData.nutrition.fat}
                                   onChange={handleNutritionChange} min="0"/>
                        </div>
                        <div className="col-md-3 mb-3"><label className="form-label">Sodium (mg)</label>
                            <input type="number" name="sodium" className="form-control"
                                   value={formData.nutrition.sodium}
                                   onChange={handleNutritionChange} min="0"/>
                        </div>
                        <div className="col-md-3 mb-3"><label className="form-label">Sugar (g)</label>
                            <input type="number" name="sugar" className="form-control" value={formData.nutrition.sugar}
                                   onChange={handleNutritionChange} min="0"/>
                        </div>
                    </div>
                </div>

                {/* --- Publish checkbox --- */}
                <div className="form-check form-switch mb-3">
                    <input
                        className="form-check-input"
                        type="checkbox"
                        role="switch"
                        id="published"
                        name="published"
                        checked={formData.published}
                        onChange={handleChange}
                    />
                    <label className="form-check-label" htmlFor="published">
                        Make Recipe Public
                    </label>
                    <small className="form-text text-muted d-block">
                        If checked, this recipe will be visible to everyone. Otherwise, only you can see it.
                    </small>
                </div>

                {/* --- Save button --- */}
                {error && <div className="alert alert-danger">{error}</div>}
                {success && <div className="alert alert-success">{success}</div>}
                <div className="d-grid">
                    <button type="submit" className="btn btn-primary btn-lg" disabled={loading}>
                        {loading ? 'Saving...' : 'Save changes'}
                    </button>
                </div>
            </form>
        </div>
    );
}