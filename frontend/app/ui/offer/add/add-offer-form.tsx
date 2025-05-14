export default function AddOfferForm() {
    return (
        <form className="flex flex-col gap-4 w-full md:w-200">
            <label htmlFor="category" className="text-lg font-semibold">Category</label>
            <select id="category" name="category" className="border rounded p-2" required defaultValue="">
                <option value="" disabled>Select a category</option>
                <option value="electronics">Electronics</option>
                <option value="furniture">Furniture</option>
                <option value="clothing">Clothing</option>
                <option value="books">Books</option>
            </select>

            <label htmlFor="description" className="text-lg font-semibold">Description</label>
            <textarea id="description" name="description" className="border rounded p-2 h-32" required></textarea>

            <button type="submit" className="bg-blue-600 text-white rounded p-2">Submit</button>
        </form>
    )
}