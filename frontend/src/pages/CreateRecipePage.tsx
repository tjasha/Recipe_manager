import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

// Ingredients from DB for the drop-down
interface Ingredient {
    id: number;
    ingredient: string;
    unit: string;
}

// Data to sent to the backend
interface RecipeFormData {
    title: string;
    description: string;
    portion: number;
    preparationTime: number;
    cookingTime: number;
    imageURL: string;
    published: boolean;
    nutrition: {
        calories: number | null;
        fat: number | null;
        sodium: number | null;
        fiber: number | null;
        carbohydrate: number | null;
        sugar: number | null;
        protein: number | null;
    };
    ingredients: {
        ingredientId: number;
        quantity: number;
    }[];
    instructions: {
        stepDescription: string;
    }[];
}

// --- Main component ---

export default function CreateRecipePage() {
    const navigate = useNavigate();

    // States

    // Main default state for the form
    const [formData, setFormData] = useState<RecipeFormData>({
        title: '',
        description: '',
        portion: 1,
        preparationTime: 0,
        cookingTime: 0,
        imageURL: '',
        published: false,
        nutrition: { calories: 0, fat: 0, sodium: 0, fiber: 0, sugar: 0, protein: 0, carbohydrate: 0 },
        ingredients: [],
        instructions: [{ stepDescription: '' }],
    });

    // State for saving all ingredients from DB
    const [allIngredients, setAllIngredients] = useState<Ingredient[]>([]);

    // State of UI
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    // --- useEffect for ingredients ---

    useEffect(() => {
        const fetchIngredients = async () => {
            try {
                const response = await fetch('http://localhost:8080/api/ingredients');
                if (!response.ok) throw new Error('Failed to fetch ingredients.');
                const data = await response.json();
                setAllIngredients(data);
            } catch (err) {
                setError('Could not load ingredients for the dropdown.');
            }
        };
        fetchIngredients();
    }, []);

    // --- Handler functions ---

    // Function to manage basic recipe fields
    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value, type } = e.target;

        if (type === 'checkbox') {
            const { checked } = e.target as HTMLInputElement;
            setFormData(prev => ({ ...prev, [name]: checked }));
            return;
        }

        let processedValue: string | number = value;
        if (type === 'number') {
            processedValue = value === '' ? 0 : parseFloat(value);
        }
        // if (name === 'portion' && processedValue < 1) {
        //     processedValue = 1;
        // }

        setFormData(prev => ({ ...prev, [name]: processedValue }));
    };

    // Function to manage nutrition object
    const handleNutritionChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            nutrition: { ...prev.nutrition, [name]: value === '' ? 0 : parseFloat(value) }
        }));
    };

    // Function to manage instructions
    const handleInstructionChange = (index: number, value: string) => {
        const newInstructions = [...formData.instructions];
        newInstructions[index].stepDescription = value;
        setFormData(prev => ({ ...prev, instructions: newInstructions }));
    };

    const addInstruction = () => {
        setFormData(prev => ({
            ...prev,
            instructions: [...prev.instructions, { stepDescription: '' }]
        }));
    };

    const removeInstruction = (index: number) => {
        setFormData(prev => ({
            ...prev,
            instructions: prev.instructions.filter((_, i) => i !== index)
        }));
    };

    // Function to manage ingredients
    const addIngredient = (ingredientId: number, quantity: number) => {
        if (!ingredientId || quantity <= 0) {
            alert('Please select an ingredient and enter a valid quantity.');
            return;
        }
        setFormData(prev => ({
            ...prev,
            ingredients: [...prev.ingredients, { ingredientId, quantity }]
        }));
    };

    const removeIngredient = (index: number) => {
        setFormData(prev => ({
            ...prev,
            ingredients: prev.ingredients.filter((_, i) => i !== index)
        }));
    };

    // function to build the form
    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (formData.title.trim() === '') {
            setError('Recipe title is required.');
            return;
        }
        if (formData.ingredients.length === 0) {
            setError('At least one ingredient is required.');
            return;
        }

        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            // TODO: create API: POST /api/recipes
            const response = await fetch('http://localhost:8080/api/createRecipe', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include', // Send cookies
                body: JSON.stringify(formData),
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || 'Failed to create recipe.');
            }

            const newRecipe = await response.json();
            setSuccess(`Recipe "${newRecipe.title}" created successfully!`);
            // Redirect to the recipe page after saving
            setTimeout(() => navigate(`/myrecipe/${newRecipe.id}`), 2000);

        } catch (err: any) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    // --- Component to add ingredients ---
    const IngredientAdder = () => {
        const [selectedId, setSelectedId] = useState<number | ''>('');
        const [quantity, setQuantity] = useState<number | ''>('');

        const handleAdd = () => {
            if (selectedId !== '' && quantity !== '' && quantity > 0) {
                addIngredient(selectedId, quantity);
                setSelectedId('');
                setQuantity('');
            }
        };

        return (
            <div className="input-group mb-3">
                <select className="form-select" value={selectedId} onChange={e => setSelectedId(Number(e.target.value))}>
                    <option value="">Choose ingredient...</option>
                    {allIngredients.map(ing => (
                        <option key={ing.id} value={ing.id}>{ing.ingredient} ({ing.unit}) </option>
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
                {/* TODO: Allows to add new ingredients */}
            </div>
        );
    };

    // --- JSX component ---
    return (
        <div className="container mt-4 mb-5">
            <h1>Create a New Recipe</h1>
            <p/>
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
                    <IngredientAdder/>
                    <ul className="list-group">
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

                {/* --- Create button --- */}
                {error && <div className="alert alert-danger">{error}</div>}
                {success && <div className="alert alert-success">{success}</div>}
                <div className="d-grid">
                    <button type="submit" className="btn btn-primary btn-lg" disabled={loading}>
                        {loading ? 'Saving...' : 'Create Recipe'}
                    </button>
                </div>
            </form>
        </div>
    );
}