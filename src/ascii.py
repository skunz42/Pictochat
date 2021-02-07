#ASCII converter

from PIL import Image
import sys

#ASCII_CHARS = ['@', '#', 'S', '%', '?', '*', '+', ';', ':', '.', ' ']
ASCII_CHARS = ['@', '#', 'S', '7', '?', '*', '+', ';', ':', '.', ' ']

def resize(image):
    (old_width, old_height) = image.size
    aspect_ratio = float(old_height)/float(old_width)
    new_height = int(aspect_ratio*80)
    new_dim = (120, new_height)
    new_image = image.resize(new_dim)
    return new_image

def grayscale(image):
    return image.convert('L')

def modify_pixels(image):
    orig_pixels = list(image.getdata())
    new_pixels = [ASCII_CHARS[pixel_value//25] for pixel_value in orig_pixels]
    return ''.join(new_pixels)

def generate_ascii(path):
    image = Image.open(path)
    image = resize(image)
    image = grayscale(image)
    pixel_arr = modify_pixels(image)
    len_pixel_arr = len(pixel_arr)
    new_image = [pixel_arr[i:i+120] for i in range(0, len_pixel_arr, 120)]
    return '\n'.join(new_image)

def print_file(text, path):
    print(path[:-4])
    f = open(path[:-4] + ".txt", "w")
    for l in text:
        f.write(l)
    f.close()

def main():
    path = sys.argv[1]
    new_image = generate_ascii(path)
    print_file(new_image, path)

main()
