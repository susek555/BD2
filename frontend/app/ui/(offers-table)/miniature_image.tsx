interface MiniatureImageProps {
    url: string;
}

export default function MainImage({ url }: MiniatureImageProps) {
    return (
        // eslint-disable-next-line @next/next/no-img-element
        <img src={url} alt="Car image" className="h-35 w-70 object-cover rounded" />
    );
}