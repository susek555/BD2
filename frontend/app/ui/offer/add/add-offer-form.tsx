export default function AddOfferForm() {
    return (
        <form className="flex flex-col gap-4 w-full md:w-200">
            <label htmlFor="title" className="text-lg font-semibold">Title</label>
            <input type="text" id="title" name="title" className="border rounded p-2" required />

            <label htmlFor="description" className="text-lg font-semibold">Description</label>
            <textarea id="description" name="description" className="border rounded p-2" required></textarea>

            <button type="submit" className="bg-blue-600 text-white rounded p-2">Submit</button>
        </form>
    )
}