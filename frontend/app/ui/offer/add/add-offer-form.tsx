"use client"

export default function AddOfferForm() {
    function handleSubmit(formData: FormData): any {
        const formDataObj = Object.fromEntries(formData.entries());
        console.log("Add Offer form data:", formDataObj);
        return formDataObj;
    }


    return (
        <form className=" w-full md:w-200" action={handleSubmit}>
            <div className="rounded-lg bg-gray-50 px-6 pb-4 pt-8 flex flex-col gap-4">
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
            </div>
        </form>
    )
}