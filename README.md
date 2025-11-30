# Stateflow

Un parser y validador para máquinas de estado (DFAs y NFAs) con soporte para funciones.

## ¿Qué es Stateflow?

Stateflow es un lenguaje de dominio específico (DSL) para definir y validar máquinas de estado deterministas (DFA) y no-deterministas (NFA). Permite:

- ✅ Definir autómatas DFA y NFA
- ✅ Declarar estados (inicial, normal, final)
- ✅ Definir transiciones con condiciones
- ✅ Crear funciones reutilizables
- ✅ Validar semánticamente el programa completo

## Uso Rápido

```bash
# Tokenizar un archivo
./stateflow tokenize example.sf

# Parsear y validar
./stateflow parse example.sf
```

## Características Principales

1. **Scanner** - Análisis léxico con soporte para strings y regex
2. **Parser** - Análisis sintáctico con validación semántica completa
3. **Symbol Table** - Tabla de símbolos con scoping anidado
4. **Validación Semántica** - 7 reglas de validación:
   - Estados iniciales únicos
   - Estados finales sin transiciones salientes
   - Determinismo en DFAs
   - Autómatas no vacíos
   - Nombres de estados únicos
   - Referencias de estados válidas
   - Parámetros de función válidos

## Ejemplo Simple

```
dfa contador {
  initial q0;
  state q1;
  final q2;
  
  on q0 -> q1 when "inc";
  on q1 -> q2 when "inc";
  on q2 -> q2 when "reset";
}

fn main(input) {
  contador <- input;
}
```