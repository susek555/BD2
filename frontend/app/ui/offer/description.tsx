export default function OfferDescription({ description }: { description: string }) {
    return (
        <div className="flex flex-col gap-4">
            <textarea
                className="text-lg resize-none overflow-auto"
                value={description}
                readOnly
                rows={15} // Zwiększono liczbę wierszy
                style={{ maxHeight: '800px' }} // Zwiększono maksymalną wysokość
            />
        </div>
    );
}