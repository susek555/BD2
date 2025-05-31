export default function OfferDescription({ description }: { description: string }) {
    const lineCount = description ? description.split('\n').length : 1;

    const estimatedRows = Math.max(
        Math.min(lineCount, 15),
        Math.ceil(description.length / 50)
    );

    return (
        <div className="flex flex-col gap-4">
            <textarea
            className="text-lg resize-none overflow-auto"
            value={description}
            readOnly
            rows={estimatedRows}
            style={{
                height: 'auto',
                minHeight: '2rem',
                maxHeight: lineCount < 7 ? '400px' : '800px',
            }}
            />
        </div>
    );
}