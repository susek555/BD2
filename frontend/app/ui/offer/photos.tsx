// in public/offer run http-server . --cors to mock urls

export default async function Photos({ imagesURLs } : { imagesURLs : string[] }) {
    return (
        <div className="flex flex-col gap-2">
            {imagesURLs.map((imageURL, index) => (
                <img
                    key={index}
                    src={imageURL}
                    alt={`Image ${index + 1}`}
                    className="w-full h-auto object-cover"
                />
            ))}
        </div>
    );
}