export default function OfferDescription({ description }: { description: string }) {
    return (
        <div className="flex flex-col gap-4">
            <textarea
            className="text-lg resize-none overflow-auto"
            value={description}
            readOnly
            rows={15}
            style={{
                maxHeight: description.split('\n').length < 7 ? '400px' : '800px', // Zmniejszono maksymalną wysokość, jeśli mniej niż 15 wierszy
            }}
            />
        </div>
    );
}