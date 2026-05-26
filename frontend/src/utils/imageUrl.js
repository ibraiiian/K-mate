/**
 * Menghasilkan URL gambar produk yang benar.
 * - Jika path sudah berupa URL lengkap (GCS), kembalikan langsung.
 * - Jika path relatif (lokal dev), tambahkan base URL server.
 * - Jika path kosong/null, kembalikan null.
 *
 * @param {string} imagePath - image_path dari API
 * @returns {string|null}
 */
export function getImageUrl(imagePath) {
    if (!imagePath) return null

    // Jika sudah berupa URL lengkap (GCS atau https/http), langsung pakai
    if (imagePath.startsWith('http://') || imagePath.startsWith('https://')) {
        return imagePath
    }

    // Path lokal relatif (dev mode): tambahkan origin saat ini
    // Di dev: origin = http://localhost:5173, tapi Vite proxy forward /uploads → backend
    // Di prod: origin = https://k-mate-app-xxx.run.app (Go serve static)
    return `${window.location.origin}${imagePath}`
}
