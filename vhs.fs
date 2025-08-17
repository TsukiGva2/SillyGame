#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

uniform sampler2D texture0;

uniform float time;

float random(vec2 st) {
    return fract(sin(dot(st.xy, vec2(12.9898,78.233))) * 43758.5453123);
}

void main() {
    vec2 texCoord = fragTexCoord;

    float jitter = sin(time * 20.0 + texCoord.y * 10.0) * 0.001;
    texCoord.x += jitter;

    float offset = 0.005; // How much the colors bleed
    float r = texture(texture0, vec2(texCoord.x + offset, texCoord.y)).r;
    float g = texture(texture0, texCoord).g;
    float b = texture(texture0, vec2(texCoord.x - offset, texCoord.y)).b;

    vec4 finalColor = vec4(r, g, b, 1.0);

    float scanline = sin(fragTexCoord.y * 800.0) * 0.1;
    finalColor.rgb -= scanline;

    float noise = (random(texCoord + time) - 0.5) * 0.15;
    finalColor.rgb += noise;

    gl_FragColor = clamp(finalColor, 0.0, 1.0);
}
